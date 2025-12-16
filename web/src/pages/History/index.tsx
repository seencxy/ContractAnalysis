import { useState } from 'react';
import {
  Card,
  Table,
  Space,
  Typography,
  Tag,
  Row,
  Col,
  Select,
  DatePicker,
  Button,
  Drawer,
  Descriptions,
  Timeline,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  SearchOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import dayjs, { type Dayjs } from 'dayjs';
import { useSignals, useSignalTracking, useSignalKlines } from '@/hooks/queries/useSignals';
import { useBinanceKlines } from '@/hooks/queries/useBinanceKlines';
import type { Signal } from '@/types/signal';
import type { KlineData } from '@/types/kline';
import { PriceLineChart } from '@/components/charts/PriceLineChart';
import { CandlestickChart } from '@/components/charts/CandlestickChart';
import StrategyBadge from '@/components/common/StrategyBadge';
import { getStatusColor, getSignalTypeColor } from '@/utils/colors';
import { useStrategies } from '@/hooks/queries/useStrategies';
import { motion } from 'framer-motion';

const { Text, Title } = Typography;
const { RangePicker } = DatePicker;

type StatusType = 'PENDING' | 'CONFIRMED' | 'TRACKING' | 'CLOSED' | 'INVALIDATED';

const statusLabels: Record<StatusType, string> = {
  PENDING: '待确认',
  CONFIRMED: '已确认',
  TRACKING: '追踪中',
  CLOSED: '已关闭',
  INVALIDATED: '已失效',
};

export default function History() {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [filters, setFilters] = useState<{
    status?: string;
    symbol?: string;
    type?: string;
    strategy?: string;
    start_time?: string;
    end_time?: string;
  }>({});

  const { data: strategies } = useStrategies();

  const [selectedSignal, setSelectedSignal] = useState<Signal | null>(null);
  const [drawerVisible, setDrawerVisible] = useState(false);
  const [klineInterval, setKlineInterval] = useState('15m');

  const { data: response, isLoading, refetch } = useSignals({ page, limit, ...filters });
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

  // Calculate start time based on interval
  const getStartTime = () => {
    if (!selectedSignal) return undefined;
    const signalTime = dayjs(selectedSignal.generated_at);
    
    let hoursBack = 24;
    switch (klineInterval) {
      case '15m': hoursBack = 50; break;
      case '1h': hoursBack = 200; break;
      case '4h': hoursBack = 800; break;
      case '1d': hoursBack = 4800; break;
      default: hoursBack = 24;
    }
    
    return signalTime.subtract(hoursBack, 'hour').valueOf();
  };

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

  const chartData: KlineData[] = (binanceKlines && binanceKlines.length > 0)
    ? binanceKlines
    : (localKlines.length > 0 && klineInterval === '15m')
        ? localKlines.map((k) => ({
            time: dayjs(k.kline_open_time).valueOf(),
            open: parseFloat(k.open_price),
            high: parseFloat(k.high_price),
            low: parseFloat(k.low_price),
            close: parseFloat(k.close_price),
            volume: parseFloat(k.volume),
          }))
        : [];

  const handleSearch = () => {
    setPage(1);
    refetch();
  };

  const handleReset = () => {
    setFilters({});
    setPage(1);
  };

  const showDetail = (record: Signal) => {
    setSelectedSignal(record);
    setKlineInterval('15m');
    setDrawerVisible(true);
  };

  const handleDateChange = (dates: [Dayjs | null, Dayjs | null] | null) => {
    if (dates && dates[0] && dates[1]) {
      setFilters({
        ...filters,
        start_time: dates[0].toISOString(),
        end_time: dates[1].toISOString(),
      });
    } else {
      const { start_time, end_time, ...rest } = filters;
      setFilters(rest);
    }
  };

  const columns: ColumnsType<Signal> = [
    {
      title: '生成时间',
      dataIndex: 'generated_at',
      key: 'generated_at',
      width: 180,
      render: (time: string) => (
        <Space direction="vertical" size={0}>
          <Text style={{ fontWeight: 500 }}>{dayjs(time).format('YYYY-MM-DD')}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {dayjs(time).format('HH:mm:ss')}
          </Text>
        </Space>
      ),
    },
    {
      title: '交易对',
      dataIndex: 'symbol',
      key: 'symbol',
      width: 120,
      render: (symbol: string) => <Text strong>{symbol}</Text>,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 90,
      align: 'center',
      render: (type: string) => (
        <Tag color={getSignalTypeColor(type)} style={{ borderRadius: 4 }}>
          {type === 'LONG' ? '做多' : '做空'}
        </Tag>
      ),
    },
    {
      title: '策略',
      dataIndex: 'strategy_name',
      key: 'strategy_name',
      width: 180,
      ellipsis: true,
      render: (_: string, record: Signal) => (
        <StrategyBadge signal={record} showDetails />
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      align: 'center',
      render: (status: StatusType) => (
        <Tag color={getStatusColor(status)} style={{ borderRadius: 4 }}>{statusLabels[status]}</Tag>
      ),
    },
    {
      title: '信号价格',
      dataIndex: 'price_at_signal',
      key: 'price_at_signal',
      width: 120,
      align: 'right',
      render: (price: string) => <Text strong>{parseFloat(price).toFixed(4)}</Text>,
    },
    {
      title: '多空比',
      key: 'tracking',
      width: 200,
      render: (_: unknown, record: Signal) => (
        <Space size="small">
          <Text type="success">{parseFloat(record.long_account_ratio).toFixed(1)}%</Text>
          <Text type="secondary">/</Text>
          <Text type="danger">{parseFloat(record.short_account_ratio).toFixed(1)}%</Text>
        </Space>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      align: 'center',
      fixed: 'right',
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
      <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 24 }}>
         <Title level={4} style={{ margin: 0 }}>历史信号查询</Title>
      </div>

      <Card bordered={false} style={{ marginBottom: 24, borderRadius: 12 }} bodyStyle={{ padding: 24 }}>
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} md={6} lg={5}>
            <RangePicker
              style={{ width: '100%' }}
              onChange={handleDateChange}
              showTime
              format="MM-DD HH:mm"
              placeholder={['开始时间', '结束时间']}
            />
          </Col>
          <Col xs={24} sm={12} md={4} lg={3}>
            <Select
              style={{ width: '100%' }}
              placeholder="交易对"
              allowClear
              value={filters.symbol}
              onChange={(value) => setFilters({ ...filters, symbol: value })}
              options={[
                { label: 'BTCUSDT', value: 'BTCUSDT' },
                { label: 'ETHUSDT', value: 'ETHUSDT' },
                { label: 'BNBUSDT', value: 'BNBUSDT' },
                { label: 'SOLUSDT', value: 'SOLUSDT' },
                { label: 'ADAUSDT', value: 'ADAUSDT' },
              ]}
            />
          </Col>
          <Col xs={24} sm={12} md={4} lg={3}>
            <Select
              style={{ width: '100%' }}
              placeholder="类型"
              allowClear
              value={filters.type}
              onChange={(value) => setFilters({ ...filters, type: value })}
              options={[
                { label: '做多', value: 'LONG' },
                { label: '做空', value: 'SHORT' },
              ]}
            />
          </Col>
          <Col xs={24} sm={12} md={4} lg={3}>
            <Select
              style={{ width: '100%' }}
              placeholder="状态"
              allowClear
              value={filters.status}
              onChange={(value) => setFilters({ ...filters, status: value })}
              options={Object.entries(statusLabels).map(([value, label]) => ({
                label,
                value,
              }))}
            />
          </Col>
          <Col xs={24} sm={24} md={6} lg={6}>
            <Select
              style={{ width: '100%' }}
              placeholder="策略"
              allowClear
              value={filters.strategy}
              onChange={(value) => setFilters({ ...filters, strategy: value })}
              options={strategies?.map(s => ({ label: s.name, value: s.key })) || []}
            />
          </Col>
          <Col xs={24} sm={24} md={24} lg={4} style={{ textAlign: 'right' }}>
            <Space>
              <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
                搜索
              </Button>
              <Button icon={<ReloadOutlined />} onClick={handleReset}>
                重置
              </Button>
            </Space>
          </Col>
        </Row>
      </Card>

      <Card bordered={false} style={{ borderRadius: 12 }} bodyStyle={{ padding: 0 }}>
        <Table
          columns={columns}
          dataSource={signals}
          rowKey="signal_id"
          loading={isLoading}
          scroll={{ x: 1200 }}
          pagination={{
            current: page,
            pageSize: limit,
            total: pagination?.total || 0,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条信号`,
            onChange: (newPage, newLimit) => {
              setPage(newPage);
              setLimit(newLimit);
            },
            pageSizeOptions: ['20', '50', '100'],
          }}
          size="middle"
        />
      </Card>

      <Drawer
        title="信号详情"
        placement="right"
        width={600}
        onClose={() => setDrawerVisible(false)}
        open={drawerVisible}
        headerStyle={{ borderBottom: 'none' }}
      >
        {selectedSignal && (
          <Space direction="vertical" size="large" style={{ width: '100%' }}>
            <Descriptions column={2} size="small">
              <Descriptions.Item label="交易对" span={2}>
                <Text strong style={{ fontSize: 16 }}>{selectedSignal.symbol}</Text>
                <Tag color={selectedSignal.type === 'LONG' ? 'green' : 'red'} style={{ marginLeft: 8 }}>
                   {selectedSignal.type === 'LONG' ? '做多' : '做空'}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="策略">
                 {selectedSignal.strategy_name}
              </Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={getStatusColor(selectedSignal.status)}>
                  {statusLabels[selectedSignal.status as StatusType]}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="信号价格">
                <Text strong>{parseFloat(selectedSignal.price_at_signal).toFixed(4)}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="生成时间">
                {dayjs(selectedSignal.generated_at).format('MM-DD HH:mm')}
              </Descriptions.Item>
            </Descriptions>

            <div style={{ background: '#fafafa', padding: 16, borderRadius: 8 }}>
                <Row gutter={[16, 16]}>
                    <Col span={12}>
                        <div style={{ fontSize: 12, color: '#999', marginBottom: 4 }}>多头账户比例</div>
                        <Text strong type="success" style={{ fontSize: 16 }}>
                            {parseFloat(selectedSignal.long_account_ratio).toFixed(2)}%
                        </Text>
                    </Col>
                    <Col span={12}>
                        <div style={{ fontSize: 12, color: '#999', marginBottom: 4 }}>空头账户比例</div>
                        <Text strong type="danger" style={{ fontSize: 16 }}>
                            {parseFloat(selectedSignal.short_account_ratio).toFixed(2)}%
                        </Text>
                    </Col>
                    <Col span={12}>
                        <div style={{ fontSize: 12, color: '#999', marginBottom: 4 }}>多头持仓比例</div>
                        <Text strong type="success" style={{ fontSize: 16 }}>
                            {parseFloat(selectedSignal.long_position_ratio).toFixed(2)}%
                        </Text>
                    </Col>
                    <Col span={12}>
                        <div style={{ fontSize: 12, color: '#999', marginBottom: 4 }}>空头持仓比例</div>
                        <Text strong type="danger" style={{ fontSize: 16 }}>
                            {parseFloat(selectedSignal.short_position_ratio).toFixed(2)}%
                        </Text>
                    </Col>
                </Row>
            </div>

            {selectedSignal.reason && (
              <div>
                <Text type="secondary" style={{ fontSize: 12, display: 'block', marginBottom: 8 }}>生成原因</Text>
                <div style={{ background: '#f5f5f5', padding: 12, borderRadius: 6, fontSize: 13, color: '#666' }}>
                  {selectedSignal.reason}
                </div>
              </div>
            )}

            {(chartData.length > 0 || trackings.length > 0) && (
              <div>
                <Text strong style={{ display: 'block', marginBottom: 16 }}>价格走势</Text>
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
              </div>
            )}

            {trackings.length > 0 && (
              <div>
                <Text strong style={{ display: 'block', marginBottom: 16 }}>追踪记录</Text>
                <Timeline
                  items={trackings.map((tracking) => ({
                    color: parseFloat(tracking.price_change_pct) >= 0 ? 'green' : 'red',
                    children: (
                      <Space direction="vertical" size={0}>
                        <Text style={{ fontSize: 12 }} type="secondary">
                          {dayjs(tracking.tracked_at).format('MM-DD HH:mm')}
                        </Text>
                        <Space>
                          <Text>{parseFloat(tracking.current_price).toFixed(4)}</Text>
                          <Text
                            strong
                            style={{
                              color: parseFloat(tracking.price_change_pct) >= 0 ? '#52c41a' : '#ff4d4f',
                            }}
                          >
                            {parseFloat(tracking.price_change_pct) >= 0 ? '+' : ''}
                            {parseFloat(tracking.price_change_pct).toFixed(2)}%
                          </Text>
                        </Space>
                      </Space>
                    ),
                  }))}
                />
              </div>
            )}
          </Space>
        )}
      </Drawer>
    </motion.div>
  );
}