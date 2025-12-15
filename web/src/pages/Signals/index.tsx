import { useState } from 'react';
import {
  Card,
  Table,
  Tag,
  Space,
  Select,
  Input,
  Button,
  Typography,
  Pagination,
  Drawer,
  Descriptions,
  Row,
  Col,
  Statistic,
  Divider,
} from 'antd';
import { ReloadOutlined, SearchOutlined, FilterOutlined, RiseOutlined, ThunderboltFilled } from '@ant-design/icons';
import { useSignals, useSignalTracking, useSignalKlines } from '@/hooks/queries/useSignals';
import { useBinanceKlines } from '@/hooks/queries/useBinanceKlines';
import type { Signal } from '@/types/signal';
import type { KlineData } from '@/types/kline';
import { Loading } from '@/components/common/Loading';
import { EmptyState } from '@/components/common/EmptyState';
import { PriceLineChart } from '@/components/charts/PriceLineChart';
import { CandlestickChart } from '@/components/charts/CandlestickChart';
import { formatTime, formatPrice, formatPercentString } from '@/utils/format';
import { getStatusColor, getSignalTypeColor } from '@/utils/colors';
import dayjs from 'dayjs';
import { motion } from 'framer-motion';

const { Text, Title } = Typography;
const { Option } = Select;

export default function Signals() {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [status, setStatus] = useState<string>('');
  const [type, setType] = useState<string>('');
  const [symbol, setSymbol] = useState<string>('');
  const [selectedSignal, setSelectedSignal] = useState<Signal | null>(null);
  const [drawerVisible, setDrawerVisible] = useState(false);
  const [klineInterval, setKlineInterval] = useState('15m');

  const { data: response, isLoading, refetch } = useSignals({
    page,
    limit,
    status: status || undefined,
    type: type || undefined,
    symbol: symbol || undefined,
  });

  const signals = response?.data?.items || [];
  const pagination = response?.data?.pagination;

  const { data: trackingRes } = useSignalTracking(selectedSignal?.signal_id || '', {
    enabled: !!selectedSignal,
  });
  const trackings = trackingRes?.data || [];

  const { data: klinesRes } = useSignalKlines(selectedSignal?.signal_id || '', {
    enabled: !!selectedSignal,
  });
  const localKlines = klinesRes?.data || [];

  // Calculate start time based on interval to ensure we get enough candles (e.g. 200)
  const getStartTime = () => {
    if (!selectedSignal) return undefined;
    const signalTime = dayjs(selectedSignal.generated_at);
    
    // Duration to look back: limit (200) * interval
    let hoursBack = 24;
    switch (klineInterval) {
      case '15m': hoursBack = 50; break; // ~200 candles
      case '1h': hoursBack = 200; break;
      case '4h': hoursBack = 800; break;
      case '1d': hoursBack = 4800; break; // ~200 days
      default: hoursBack = 24;
    }
    
    return signalTime.subtract(hoursBack, 'hour').valueOf();
  };

  // Fetch Binance Klines
  const { data: binanceKlines } = useBinanceKlines(
    {
      symbol: selectedSignal?.symbol || '',
      interval: klineInterval,
      startTime: getStartTime(),
      limit: 500,
    },
    {
      enabled: !!selectedSignal,
    }
  );

  // Prefer Binance klines for dynamic interval support, fallback to local if binance fails/empty
  const chartData: KlineData[] = (binanceKlines && binanceKlines.length > 0)
    ? binanceKlines
    : (localKlines.length > 0 && klineInterval === '15m') // Only use local if interval matches (assuming local is 15m)
        ? localKlines.map((k) => ({
            time: dayjs(k.kline_open_time).valueOf(),
            open: parseFloat(k.open_price),
            high: parseFloat(k.high_price),
            low: parseFloat(k.low_price),
            close: parseFloat(k.close_price),
            volume: parseFloat(k.volume),
          }))
        : [];

  const showDetail = (record: Signal) => {
    setSelectedSignal(record);
    setKlineInterval('15m'); // Reset to default when opening details
    setDrawerVisible(true);
  };

  const columns = [
    {
      title: '交易对',
      dataIndex: 'symbol',
      key: 'symbol',
      render: (text: string) => <Text strong style={{ fontSize: 15 }}>{text}</Text>,
      width: 120,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color={getSignalTypeColor(type)} style={{ borderRadius: 6, padding: '0 8px', fontWeight: 500 }}>{type}</Tag>
      ),
      width: 90,
    },
    {
      title: '策略',
      dataIndex: 'strategy_name',
      key: 'strategy_name',
      width: 180,
      render: (text: string) => <Text type="secondary" style={{ fontSize: 13 }}>{text}</Text>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={getStatusColor(status)} style={{ borderRadius: 6 }}>{status}</Tag>
      ),
      width: 110,
    },
    {
      title: '信号价格',
      dataIndex: 'price_at_signal',
      key: 'price_at_signal',
      render: (price: string) => <Text strong>{formatPrice(price)}</Text>,
      width: 120,
    },
    {
      title: '多账户比',
      dataIndex: 'long_account_ratio',
      key: 'long_account_ratio',
      render: (ratio: string) => <Text type="success">{formatPercentString(ratio)}</Text>,
      width: 100,
    },
    {
      title: '空账户比',
      dataIndex: 'short_account_ratio',
      key: 'short_account_ratio',
      render: (ratio: string) => <Text type="danger">{formatPercentString(ratio)}</Text>,
      width: 100,
    },
    {
      title: '生成时间',
      dataIndex: 'generated_at',
      key: 'generated_at',
      render: (time: string) => <Text type="secondary" style={{ fontSize: 13 }}>{formatTime(time)}</Text>,
      width: 180,
    },
    {
      title: '确认状态',
      dataIndex: 'is_confirmed',
      key: 'is_confirmed',
      render: (confirmed: boolean) => (
        <Tag bordered={false} color={confirmed ? 'success' : 'warning'} style={{ borderRadius: 4 }}>
          {confirmed ? '已确认' : '待确认'}
        </Tag>
      ),
      width: 100,
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      fixed: 'right' as const,
      render: (_: unknown, record: Signal) => (
        <Button type="link" size="small" onClick={() => showDetail(record)}>
          详情
        </Button>
      ),
    },
  ];

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4 }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={4} style={{ margin: 0 }}>实时信号监控</Title>
        <Button
          icon={<ReloadOutlined />}
          onClick={() => refetch()}
          loading={isLoading}
        >
          刷新数据
        </Button>
      </div>

      <Card bordered={false} style={{ marginBottom: 24, borderRadius: 12 }} bodyStyle={{ padding: 24 }}>
        <Space wrap size="middle">
            <Input
                prefix={<SearchOutlined style={{ color: '#bfbfbf' }} />}
                placeholder="搜索交易对..."
                value={symbol}
                onChange={(e) => setSymbol(e.target.value)}
                allowClear
                style={{ width: 240 }}
            />
          <Select
            style={{ width: 160 }}
            placeholder="选择状态"
            value={status || undefined}
            onChange={setStatus}
            allowClear
            suffixIcon={<FilterOutlined style={{ color: '#bfbfbf' }} />}
          >
            <Option value="">全部状态</Option>
            <Option value="PENDING">待确认</Option>
            <Option value="CONFIRMED">已确认</Option>
            <Option value="TRACKING">追踪中</Option>
            <Option value="CLOSED">已关闭</Option>
            <Option value="INVALIDATED">已失效</Option>
          </Select>

          <Select
            style={{ width: 160 }}
            placeholder="选择类型"
            value={type || undefined}
            onChange={setType}
            allowClear
          >
            <Option value="">全部类型</Option>
            <Option value="LONG">做多</Option>
            <Option value="SHORT">做空</Option>
          </Select>
        </Space>
      </Card>

      <Card bordered={false} style={{ borderRadius: 12, overflow: 'hidden' }} bodyStyle={{ padding: 0 }}>
        {isLoading ? (
          <div style={{ padding: 48 }}>
            <Loading />
          </div>
        ) : signals.length === 0 ? (
          <EmptyState message="暂无信号数据" />
        ) : (
          <>
            <Table
              dataSource={signals}
              columns={columns}
              rowKey="signal_id"
              pagination={false}
              scroll={{ x: 1300 }}
              size="middle"
            />
            {pagination && (
              <div style={{ padding: '24px', display: 'flex', justifyContent: 'flex-end' }}>
                <Pagination
                  current={pagination.page}
                  pageSize={pagination.limit}
                  total={pagination.total}
                  onChange={(newPage, newPageSize) => {
                    setPage(newPage);
                    if (newPageSize) setLimit(newPageSize);
                  }}
                  showSizeChanger
                  showTotal={(total) => `共 ${total} 条`}
                />
              </div>
            )}
          </>
        )}
      </Card>

      <Drawer
        title={<span style={{ fontSize: 18, fontWeight: 600 }}>信号详情</span>}
        placement="right"
        width={720}
        onClose={() => setDrawerVisible(false)}
        open={drawerVisible}
        styles={{ 
            header: { borderBottom: '1px solid #f0f0f0', padding: '20px 24px' },
            body: { padding: '24px' }
        }}
      >
        {selectedSignal && (
          <Space direction="vertical" size="large" style={{ width: '100%' }}>
            
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                <div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 8 }}>
                        <Title level={3} style={{ margin: 0 }}>{selectedSignal.symbol}</Title>
                        <Tag color={getSignalTypeColor(selectedSignal.type)} style={{ fontSize: 14, padding: '4px 10px' }}>
                            {selectedSignal.type === 'LONG' ? '做多' : '做空'}
                        </Tag>
                    </div>
                    <Text type="secondary">ID: {selectedSignal.signal_id}</Text>
                </div>
                <div style={{ textAlign: 'right' }}>
                    <div style={{ fontSize: 24, fontWeight: 700, color: '#1677ff' }}>
                        {formatPrice(selectedSignal.price_at_signal)}
                    </div>
                    <Text type="secondary">信号价格</Text>
                </div>
            </div>

            <Divider style={{ margin: '12px 0' }} />

            <Row gutter={[24, 24]}>
                <Col span={12}>
                    <Descriptions title="基础信息" column={1} size="small">
                        <Descriptions.Item label="策略名称">{selectedSignal.strategy_name}</Descriptions.Item>
                        <Descriptions.Item label="当前状态">
                            <Tag color={getStatusColor(selectedSignal.status)}>{selectedSignal.status}</Tag>
                        </Descriptions.Item>
                        <Descriptions.Item label="生成时间">{formatTime(selectedSignal.generated_at)}</Descriptions.Item>
                    </Descriptions>
                </Col>
                <Col span={12}>
                    <Descriptions title="市场情绪" column={1} size="small">
                        <Descriptions.Item label="多头账户比">
                            <Text type="success">{formatPercentString(selectedSignal.long_account_ratio)}</Text>
                        </Descriptions.Item>
                        <Descriptions.Item label="空头账户比">
                            <Text type="danger">{formatPercentString(selectedSignal.short_account_ratio)}</Text>
                        </Descriptions.Item>
                        <Descriptions.Item label="多空持仓比">
                            <Text type="success">{formatPercentString(selectedSignal.long_position_ratio)}</Text> / <Text type="danger">{formatPercentString(selectedSignal.short_position_ratio)}</Text>
                        </Descriptions.Item>
                    </Descriptions>
                </Col>
            </Row>

            {selectedSignal.reason && (
              <div style={{ background: '#f5f7fa', padding: 16, borderRadius: 8 }}>
                <Text strong style={{ display: 'block', marginBottom: 8 }}>生成原因</Text>
                <Text style={{ color: '#595959', fontSize: 13 }}>{selectedSignal.reason}</Text>
              </div>
            )}

            {(chartData.length > 0 || trackings.length > 0) && (
              <Card 
                title="价格走势" 
                size="small" 
                bordered={false} 
                style={{ background: '#fff', boxShadow: '0 2px 8px rgba(0,0,0,0.04)' }}
              >
                {chartData.length > 0 ? (
                  <CandlestickChart
                    klines={chartData}
                    signalPrice={selectedSignal.price_at_signal}
                    signalType={selectedSignal.type}
                    signalTime={selectedSignal.generated_at}
                    interval={klineInterval}
                    onIntervalChange={setKlineInterval}
                  />
                ) : (
                  <PriceLineChart
                    trackings={trackings}
                    signalPrice={selectedSignal.price_at_signal}
                    signalType={selectedSignal.type}
                  />
                )}
              </Card>
            )}

            {trackings.length > 0 && (
              <div style={{ marginTop: 16 }}>
                <Title level={5} style={{ marginBottom: 16 }}>追踪统计</Title>
                <Row gutter={16}>
                  <Col span={8}>
                    <Statistic 
                        title="最新价格" 
                        value={formatPrice(trackings[trackings.length - 1].current_price)} 
                        valueStyle={{ fontSize: 18, fontWeight: 600 }}
                    />
                  </Col>
                  <Col span={8}>
                    <Statistic 
                        title="价格变化" 
                        value={formatPercentString(trackings[trackings.length - 1].price_change_pct)}
                        valueStyle={{ color: parseFloat(trackings[trackings.length - 1].price_change_pct) >= 0 ? '#52c41a' : '#f5222d', fontSize: 18, fontWeight: 600 }}
                        prefix={parseFloat(trackings[trackings.length - 1].price_change_pct) >= 0 ? <RiseOutlined /> : <ThunderboltFilled />}
                    />
                  </Col>
                  <Col span={8}>
                    <Statistic 
                        title="追踪时长 (小时)" 
                        value={trackings[trackings.length - 1].hours_tracked} 
                        valueStyle={{ fontSize: 18 }}
                    />
                  </Col>
                </Row>
              </div>
            )}
          </Space>
        )}
      </Drawer>
    </motion.div>
  );
}