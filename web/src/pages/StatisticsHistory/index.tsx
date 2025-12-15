import { useState, useMemo } from 'react';
import {
  Card,
  DatePicker,
  Select,
  Table,
  Space,
  Typography,
  Tabs,
  Tag,
  Row,
  Col,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import dayjs, { Dayjs } from 'dayjs';
import { useStatisticsHistory } from '@/hooks/queries/useStatistics';
import { StatisticsTrendChart } from '@/components/charts/StatisticsTrendChart';
import type { Statistics } from '@/types/statistics';
import { formatPercentString } from '@/utils/format';
import { Loading } from '@/components/common/Loading';
import { EmptyState } from '@/components/common/EmptyState';
import { motion } from 'framer-motion';

const { RangePicker } = DatePicker;
const { Title, Text } = Typography;
const { Option } = Select;

export default function StatisticsHistory() {
  // Filter states
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs]>([
    dayjs().subtract(7, 'days'),
    dayjs(),
  ]);
  const [strategy, setStrategy] = useState<string | undefined>();
  const [symbol, setSymbol] = useState<string | undefined>();

  // Fetch data
  const { data: response, isLoading } = useStatisticsHistory({
    start_time: dateRange[0].toISOString(),
    end_time: dateRange[1].toISOString(),
    strategy,
    symbol,
  });

  const historyData = response?.data || [];

  // Extract unique strategies and symbols for filters
  const { strategies, symbols } = useMemo(() => {
    const strategySet = new Set<string>();
    const symbolSet = new Set<string>();

    historyData.forEach((item) => {
      strategySet.add(item.strategy_name);
      if (item.symbol) symbolSet.add(item.symbol);
    });

    return {
      strategies: Array.from(strategySet),
      symbols: Array.from(symbolSet),
    };
  }, [historyData]);

  // Table columns
  const columns: ColumnsType<Statistics> = [
    {
      title: '计算时间',
      dataIndex: 'calculated_at',
      key: 'calculated_at',
      width: 160,
      render: (text: string) => <Text>{dayjs(text).format('YYYY-MM-DD HH:mm')}</Text>,
      sorter: (a, b) => dayjs(a.calculated_at).unix() - dayjs(b.calculated_at).unix(),
      defaultSortOrder: 'descend',
    },
    {
      title: '策略',
      dataIndex: 'strategy_name',
      key: 'strategy_name',
      width: 150,
      render: (text) => <Text type="secondary">{text}</Text>,
    },
    {
      title: '交易对',
      dataIndex: 'symbol',
      key: 'symbol',
      width: 120,
      render: (text: string) => text ? <Text strong>{text}</Text> : <Tag>全部</Tag>,
    },
    {
      title: '总信号',
      dataIndex: 'total_signals',
      key: 'total_signals',
      width: 100,
      align: 'right',
    },
    {
      title: '胜率',
      dataIndex: 'win_rate',
      key: 'win_rate',
      width: 120,
      align: 'center',
      render: (text: string) => {
        if (!text) return '-';
        const rate = parseFloat(text);
        let color = 'default';
        if (rate >= 50) color = 'success';
        else if (rate >= 30) color = 'warning';
        else color = 'error';
        return <Tag color={color} bordered={false}>{formatPercentString(text)}</Tag>;
      },
      sorter: (a, b) => parseFloat(a.win_rate || '0') - parseFloat(b.win_rate || '0'),
    },
    {
      title: '盈/亏',
      key: 'pnl_count',
      width: 140,
      align: 'center',
      render: (_: unknown, record: Statistics) => (
        <Space size={4}>
            <span style={{ color: '#52c41a' }}>{record.profitable_signals}</span>
            <span style={{ color: '#bfbfbf' }}>/</span>
            <span style={{ color: '#ff4d4f' }}>{record.losing_signals}</span>
        </Space>
      )
    },
    {
      title: '平均盈利',
      dataIndex: 'avg_profit_pct',
      key: 'avg_profit_pct',
      width: 120,
      align: 'right',
      render: (text: string) => (text ? <Text type="success">{formatPercentString(text)}</Text> : '-'),
    },
    {
      title: '盈利因子',
      dataIndex: 'profit_factor',
      key: 'profit_factor',
      width: 100,
      align: 'right',
      render: (text: string) => (text ? parseFloat(text).toFixed(2) : '-'),
    },
  ];

  return (
    <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4 }}
    >
      <div style={{ marginBottom: 24 }}>
        <Title level={4} style={{ margin: 0 }}>历史统计分析</Title>
      </div>

      <Card bordered={false} style={{ marginBottom: 24, borderRadius: 12 }} bodyStyle={{ padding: 24 }}>
        <Row gutter={[24, 24]} align="middle">
          <Col xs={24} md={8}>
            <div style={{ marginBottom: 8, fontSize: 12, color: '#8c8c8c' }}>时间范围</div>
            <RangePicker
              value={dateRange}
              onChange={(dates) => dates && setDateRange(dates as [Dayjs, Dayjs])}
              showTime
              format="YYYY-MM-DD HH:mm"
              style={{ width: '100%' }}
            />
          </Col>

          <Col xs={24} md={8}>
            <div style={{ marginBottom: 8, fontSize: 12, color: '#8c8c8c' }}>策略选择</div>
            <Select
              style={{ width: '100%' }}
              placeholder="全部策略"
              allowClear
              value={strategy}
              onChange={setStrategy}
            >
              {strategies.map((s) => (
                <Option key={s} value={s}>
                  {s}
                </Option>
              ))}
            </Select>
          </Col>

          <Col xs={24} md={8}>
            <div style={{ marginBottom: 8, fontSize: 12, color: '#8c8c8c' }}>交易对</div>
            <Select
              style={{ width: '100%' }}
              placeholder="全部交易对"
              allowClear
              value={symbol}
              onChange={setSymbol}
            >
              {symbols.map((s) => (
                <Option key={s} value={s}>
                  {s}
                </Option>
              ))}
            </Select>
          </Col>
        </Row>
      </Card>

      {/* Charts */}
      {isLoading ? (
        <Loading />
      ) : historyData.length === 0 ? (
        <EmptyState message="暂无历史统计数据" />
      ) : (
        <>
          <Card 
            title="趋势分析" 
            bordered={false} 
            style={{ marginBottom: 24, borderRadius: 12 }}
          >
            <Tabs
              items={[
                {
                  key: 'win_rate',
                  label: '胜率趋势',
                  children: <StatisticsTrendChart data={historyData} metricType="win_rate" />,
                },
                {
                  key: 'total_signals',
                  label: '信号数量',
                  children: <StatisticsTrendChart data={historyData} metricType="total_signals" />,
                },
                {
                  key: 'avg_profit',
                  label: '平均盈利',
                  children: <StatisticsTrendChart data={historyData} metricType="avg_profit" />,
                },
                {
                  key: 'profit_factor',
                  label: '盈利因子',
                  children: <StatisticsTrendChart data={historyData} metricType="profit_factor" />,
                },
              ]}
              tabBarStyle={{ marginBottom: 24 }}
            />
          </Card>

          {/* Data Table */}
          <Card 
            title="详细记录" 
            bordered={false} 
            style={{ borderRadius: 12 }}
            bodyStyle={{ padding: 0 }}
          >
            <Table
              columns={columns}
              dataSource={historyData}
              rowKey={(record) =>
                `${record.strategy_name}-${record.symbol || 'all'}-${record.period_label}-${record.calculated_at}`
              }
              pagination={{
                pageSize: 15,
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 条记录`,
              }}
              loading={isLoading}
              scroll={{ x: 1200 }}
              size="middle"
            />
          </Card>
        </>
      )}
    </motion.div>
  );
}