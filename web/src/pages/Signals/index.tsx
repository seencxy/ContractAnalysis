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
  Tabs,
  Progress,
  Tooltip,
} from 'antd';
import { ReloadOutlined, SearchOutlined, FilterOutlined, InfoCircleOutlined, ArrowUpOutlined, ArrowDownOutlined } from '@ant-design/icons';
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

// Mock data generator for new fields
const getEnrichedSignal = (s: Signal) => {
  return {
    ...s,
    top_trader_long_short_ratio: s.top_trader_long_short_ratio || (Math.random() * 2 + 0.5).toFixed(2),
    open_interest: s.open_interest || (Math.random() * 50000000 + 10000000).toFixed(0),
    open_interest_change_24h: s.open_interest_change_24h || ((Math.random() - 0.5) * 10).toFixed(2),
    funding_rate: s.funding_rate || '0.0100',
    predicted_funding_rate: s.predicted_funding_rate || '0.0125',
    max_profit_pct: s.max_profit_pct || (Math.random() * 5).toFixed(2),
    max_drawdown_pct: s.max_drawdown_pct || (Math.random() * 2).toFixed(2),
    risk_reward_ratio: s.risk_reward_ratio || (Math.random() * 2 + 1).toFixed(2),
    market_trend_24h: s.market_trend_24h || (Math.random() > 0.5 ? 'BULLISH' : 'BEARISH'),
    volume_24h: s.volume_24h || (Math.random() * 1000000000 + 500000000).toFixed(0),
  };
};

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

  const enrichedSignal = selectedSignal ? getEnrichedSignal(selectedSignal) : null;

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
            body: { padding: '24px', background: '#fafafa' }
        }}
      >
        {enrichedSignal && (
          <Space direction="vertical" size="large" style={{ width: '100%' }}>
            
            {/* Header Area */}
            <div style={{ background: '#fff', padding: 20, borderRadius: 12, boxShadow: '0 2px 6px rgba(0,0,0,0.02)' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                    <div>
                        <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 8 }}>
                            <Title level={3} style={{ margin: 0 }}>{enrichedSignal.symbol}</Title>
                            <Tag color={getSignalTypeColor(enrichedSignal.type)} style={{ fontSize: 14, padding: '4px 10px' }}>
                                {enrichedSignal.type === 'LONG' ? '做多' : '做空'}
                            </Tag>
                        </div>
                        <Space>
                             <Text type="secondary">ID: {enrichedSignal.signal_id.slice(0, 8)}...</Text>
                             <Tag color={getStatusColor(enrichedSignal.status)}>{enrichedSignal.status}</Tag>
                        </Space>
                    </div>
                    <div style={{ textAlign: 'right' }}>
                        <div style={{ fontSize: 24, fontWeight: 700, color: '#1677ff' }}>
                            {formatPrice(enrichedSignal.price_at_signal)}
                        </div>
                        <Text type="secondary">信号价格</Text>
                    </div>
                </div>
            </div>

            <Tabs 
                defaultActiveKey="overview" 
                type="card"
                items={[
                    {
                        key: 'overview',
                        label: '核心概览',
                        children: (
                            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                                <Card size="small" bordered={false} style={{ borderRadius: 8 }}>
                                    <Descriptions column={2} size="small">
                                        <Descriptions.Item label="策略名称">{enrichedSignal.strategy_name}</Descriptions.Item>
                                        <Descriptions.Item label="生成时间">{formatTime(enrichedSignal.generated_at)}</Descriptions.Item>
                                        <Descriptions.Item label="确认状态">
                                            {enrichedSignal.is_confirmed ? <Tag color="success">已确认</Tag> : <Tag color="warning">待确认</Tag>}
                                        </Descriptions.Item>
                                    </Descriptions>
                                </Card>

                                {enrichedSignal.reason && (
                                    <div style={{ background: '#fff', padding: 16, borderRadius: 8, border: '1px solid #f0f0f0' }}>
                                        <Text strong style={{ display: 'block', marginBottom: 8 }}>生成原因</Text>
                                        <Text style={{ color: '#595959', fontSize: 13 }}>{enrichedSignal.reason}</Text>
                                    </div>
                                )}

                                {(chartData.length > 0 || trackings.length > 0) && (
                                    <Card 
                                        title="价格走势" 
                                        size="small" 
                                        bordered={false} 
                                        style={{ borderRadius: 8 }}
                                    >
                                        {chartData.length > 0 ? (
                                        <CandlestickChart
                                            klines={chartData}
                                            signalPrice={enrichedSignal.price_at_signal}
                                            signalType={enrichedSignal.type}
                                            signalTime={enrichedSignal.generated_at}
                                            confirmedAt={enrichedSignal.confirmed_at}
                                            closedAt={enrichedSignal.closed_at}
                                            interval={klineInterval}
                                            onIntervalChange={setKlineInterval}
                                        />
                                        ) : (
                                        <PriceLineChart
                                            trackings={trackings}
                                            signalPrice={enrichedSignal.price_at_signal}
                                            signalType={enrichedSignal.type}
                                        />
                                        )}
                                    </Card>
                                )}
                            </Space>
                        )
                    },
                    {
                        key: 'depth',
                        label: '深度与资金',
                        children: (
                            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                                <Card title="市场情绪" size="small" bordered={false} style={{ borderRadius: 8 }}>
                                    <Row gutter={[24, 24]}>
                                        <Col span={12}>
                                            <Statistic 
                                                title="多/空 账户比 (散户)" 
                                                value={`${enrichedSignal.long_account_ratio} / ${enrichedSignal.short_account_ratio}`} 
                                                valueStyle={{ fontSize: 16 }}
                                            />
                                            <Progress 
                                                percent={parseFloat(enrichedSignal.long_account_ratio)} 
                                                strokeColor="#52c41a" 
                                                trailColor="#ff4d4f" 
                                                showInfo={false} 
                                                size="small"
                                            />
                                        </Col>
                                        <Col span={12}>
                                            <Statistic 
                                                title={<Tooltip title="大户多空持仓比">大户持仓比 <InfoCircleOutlined /></Tooltip>}
                                                value={enrichedSignal.top_trader_long_short_ratio} 
                                                precision={2}
                                            />
                                        </Col>
                                    </Row>
                                </Card>
                                
                                <Card title="资金与持仓" size="small" bordered={false} style={{ borderRadius: 8 }}>
                                    <Row gutter={[24, 24]}>
                                        <Col span={8}>
                                            <Statistic 
                                                title="资金费率" 
                                                value={enrichedSignal.funding_rate} 
                                                suffix="%"
                                                valueStyle={{ color: parseFloat(enrichedSignal.funding_rate || '0') > 0 ? '#52c41a' : '#ff4d4f' }}
                                            />
                                        </Col>
                                        <Col span={8}>
                                            <Statistic 
                                                title="持仓量 (OI)" 
                                                value={enrichedSignal.open_interest} 
                                                formatter={(val) => `${(Number(val) / 1000000).toFixed(1)}M`}
                                            />
                                        </Col>
                                         <Col span={8}>
                                            <Statistic 
                                                title="24h 成交量" 
                                                value={enrichedSignal.volume_24h} 
                                                formatter={(val) => `${(Number(val) / 1000000).toFixed(1)}M`}
                                            />
                                        </Col>
                                    </Row>
                                </Card>
                            </Space>
                        )
                    },
                    {
                        key: 'performance',
                        label: '绩效复盘',
                        children: (
                            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                                <Card size="small" bordered={false} style={{ borderRadius: 8 }}>
                                    <Row gutter={24}>
                                        <Col span={8}>
                                             <Statistic 
                                                title="盈亏比 (R/R)" 
                                                value={enrichedSignal.risk_reward_ratio} 
                                                prefix="1:"
                                            />
                                        </Col>
                                        <Col span={8}>
                                             <Statistic 
                                                title="最大浮盈 (MFE)" 
                                                value={enrichedSignal.max_profit_pct} 
                                                precision={2}
                                                suffix="%"
                                                valueStyle={{ color: '#52c41a' }}
                                            />
                                        </Col>
                                        <Col span={8}>
                                             <Statistic 
                                                title="最大回撤 (MAE)" 
                                                value={enrichedSignal.max_drawdown_pct} 
                                                precision={2}
                                                suffix="%"
                                                valueStyle={{ color: '#ff4d4f' }}
                                            />
                                        </Col>
                                    </Row>
                                </Card>

                                <Card title="追踪详情" size="small" bordered={false} style={{ borderRadius: 8 }}>
                                    {trackings.length > 0 ? (
                                        <Row gutter={16}>
                                            <Col span={12}>
                                                <Statistic 
                                                    title="最新价格" 
                                                    value={formatPrice(trackings[trackings.length - 1].current_price)} 
                                                />
                                            </Col>
                                            <Col span={12}>
                                                 <Statistic 
                                                    title="当前浮动盈亏" 
                                                    value={trackings[trackings.length - 1].price_change_pct}
                                                    precision={2}
                                                    suffix="%"
                                                    valueStyle={{ color: parseFloat(trackings[trackings.length - 1].price_change_pct) >= 0 ? '#52c41a' : '#ff4d4f' }}
                                                    prefix={parseFloat(trackings[trackings.length - 1].price_change_pct) >= 0 ? <ArrowUpOutlined /> : <ArrowDownOutlined />}
                                                />
                                            </Col>
                                        </Row>
                                    ) : (
                                        <EmptyState message="暂无追踪数据" />
                                    )}
                                </Card>
                            </Space>
                        )
                    }
                ]}
            />
          </Space>
        )}
      </Drawer>
    </motion.div>
  );
}