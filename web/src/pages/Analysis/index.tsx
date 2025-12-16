import { useState, useMemo } from 'react';
import { Card, Radio, Table, Row, Col, Space, Typography, Tag, Statistic, Select } from 'antd';
import type { RadioChangeEvent } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useStrategyStatistics } from '@/hooks/queries/useStatistics';
import { useStrategies } from '@/hooks/queries/useStrategies';
import type { Statistics } from '@/types/statistics';
import { Loading } from '@/components/common/Loading';
import { EmptyState } from '@/components/common/EmptyState';
import dayjs from 'dayjs';
import { motion } from 'framer-motion';

const { Text, Title } = Typography;
const { Option } = Select;

type PeriodType = '24h' | '7d' | '30d' | 'all';

export default function Analysis() {
  const [period, setPeriod] = useState<PeriodType>('24h');
  const [strategy, setStrategy] = useState<string | undefined>(undefined);

  const { data: strategies } = useStrategies();

  const { data: response, isLoading } = useStrategyStatistics({ period, strategy });
  const allData = response?.data || [];

  // Separate summary (no symbol) from per-symbol data
  const { summary, symbolStats } = useMemo(() => {
    const summaryItem = allData.find((item) => !item.symbol);
    const symbols = allData.filter((item) => !!item.symbol);
    return { summary: summaryItem, symbolStats: symbols };
  }, [allData]);

  const handlePeriodChange = (e: RadioChangeEvent) => {
    setPeriod(e.target.value);
  };

  const columns: ColumnsType<Statistics> = [
    {
      title: '交易对',
      dataIndex: 'symbol',
      key: 'symbol',
      fixed: 'left',
      width: 140,
      render: (text: string) => <Text strong>{text}</Text>,
    },
    {
      title: '总信号',
      dataIndex: 'total_signals',
      key: 'total_signals',
      align: 'center',
      width: 100,
      sorter: (a, b) => a.total_signals - b.total_signals,
    },
    {
      title: '已确认',
      dataIndex: 'confirmed_signals',
      key: 'confirmed_signals',
      align: 'center',
      width: 120,
      render: (val: number, record: Statistics) => (
        <Space direction="vertical" size={0}>
          <Text type={val > 0 ? 'success' : 'secondary'} strong>{val}</Text>
          <Text type="secondary" style={{ fontSize: 11 }}>
            {record.total_signals > 0 ? ((val / record.total_signals) * 100).toFixed(0) : 0}%
          </Text>
        </Space>
      ),
    },
    {
      title: '盈/亏/平',
      key: 'outcomes',
      align: 'center',
      width: 150,
      render: (_: unknown, record: Statistics) => (
        <Space size={8}>
          <Tag color="success" bordered={false}>{record.profitable_signals}</Tag>
          <Tag color="error" bordered={false}>{record.losing_signals}</Tag>
          <Tag color="default" bordered={false}>{record.neutral_signals}</Tag>
        </Space>
      ),
    },
    {
      title: '胜率',
      dataIndex: 'win_rate',
      key: 'win_rate',
      align: 'center',
      width: 110,
      render: (value?: string) => {
        if (!value || value === '0') return <Text type="secondary">-</Text>;
        const rate = parseFloat(value);
        return <Tag color={rate >= 50 ? 'success' : rate >= 30 ? 'warning' : 'error'}>{rate.toFixed(1)}%</Tag>;
      },
      sorter: (a, b) => parseFloat(a.win_rate || '0') - parseFloat(b.win_rate || '0'),
    },
    {
      title: '平均持仓',
      dataIndex: 'avg_holding_hours',
      key: 'avg_holding_hours',
      align: 'center',
      width: 120,
      render: (value?: string) => (value ? `${parseFloat(value).toFixed(1)}h` : '-'),
      sorter: (a, b) => parseFloat(a.avg_holding_hours || '0') - parseFloat(b.avg_holding_hours || '0'),
    },
    {
      title: '盈利K线',
      key: 'profitable_klines',
      align: 'center',
      width: 140,
      render: (_: unknown, record: Statistics) => (
        <Space direction="vertical" size={0}>
          <Text style={{ fontSize: 12 }}>High: <Text type="success">{record.profitable_kline_hours_high}</Text></Text>
          <Text style={{ fontSize: 12 }}>Close: <Text>{record.profitable_kline_hours_close}</Text></Text>
        </Space>
      ),
    },
  ];

  if (isLoading) {
    return <Loading />;
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4 }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={4} style={{ margin: 0 }}>策略效果分析</Title>
      </div>

      <Card bordered={false} style={{ marginBottom: 24, borderRadius: 12 }} bodyStyle={{ padding: '16px 24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 24, flexWrap: 'wrap' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <Text strong>统计周期:</Text>
            <Radio.Group value={period} onChange={handlePeriodChange} buttonStyle="solid">
              <Radio.Button value="24h">24小时</Radio.Button>
              <Radio.Button value="7d">7天</Radio.Button>
              <Radio.Button value="30d">30天</Radio.Button>
              <Radio.Button value="all">全部</Radio.Button>
            </Radio.Group>
          </div>

          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <Text strong>策略筛选:</Text>
            <Select
              style={{ width: 200 }}
              placeholder="全部策略"
              allowClear
              value={strategy}
              onChange={setStrategy}
            >
              <Option value={undefined}>全部策略</Option>
              {strategies?.map(s => <Option key={s.key} value={s.key}>{s.name}</Option>)}
            </Select>
          </div>
        </div>
      </Card>

      {summary && (
        <Card 
          title={
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Text strong style={{ fontSize: 16 }}>{summary.strategy_name || '总体统计'}</Text>
              <Tag color="blue">{dayjs(summary.period_start).format('MM-DD HH:mm')} ~ {dayjs(summary.period_end).format('MM-DD HH:mm')}</Tag>
            </div>
          }
          bordered={false}
          style={{ marginBottom: 24, borderRadius: 12 }}
        >
          <Row gutter={[24, 24]}>
            <Col xs={12} sm={8} md={4}>
              <Statistic title="总信号" value={summary.total_signals} />
            </Col>
            <Col xs={12} sm={8} md={4}>
              <Statistic 
                title="已确认" 
                value={summary.confirmed_signals} 
                suffix={<span style={{ fontSize: 12, color: '#999' }}>/ {summary.total_signals}</span>}
                valueStyle={{ color: '#52c41a' }}
              />
            </Col>
            <Col xs={12} sm={8} md={4}>
              <Statistic 
                title="胜率" 
                value={summary.win_rate ? parseFloat(summary.win_rate).toFixed(1) : 0} 
                suffix="%" 
                valueStyle={{ color: parseFloat(summary.win_rate || '0') >= 50 ? '#52c41a' : '#faad14' }}
              />
            </Col>
            <Col xs={12} sm={8} md={4}>
              <Statistic 
                title="盈利/亏损" 
                value={summary.profitable_signals}
                suffix={<span style={{ fontSize: 14, color: '#ff4d4f' }}>/ {summary.losing_signals}</span>}
                valueStyle={{ color: '#52c41a' }}
              />
            </Col>
             <Col xs={12} sm={8} md={4}>
              <Statistic 
                title="中性" 
                value={summary.neutral_signals}
                valueStyle={{ color: '#faad14' }}
              />
            </Col>
            <Col xs={12} sm={8} md={4}>
              <Statistic 
                title="平均持仓" 
                value={summary.avg_holding_hours ? parseFloat(summary.avg_holding_hours).toFixed(1) : '-'} 
                suffix="h"
              />
            </Col>
          </Row>
        </Card>
      )}

      <Card 
        title={`交易对详细统计 (${symbolStats.length})`} 
        bordered={false} 
        style={{ borderRadius: 12 }}
        bodyStyle={{ padding: 0 }}
      >
        {symbolStats.length === 0 ? (
          <EmptyState message="暂无数据" />
        ) : (
          <Table
            columns={columns}
            dataSource={symbolStats}
            rowKey={(record) => `${record.symbol}-${record.period_label}`}
            scroll={{ x: 1000 }}
            pagination={{
              pageSize: 15,
              showSizeChanger: true,
              showTotal: (total) => `共 ${total} 个交易对`,
            }}
            size="middle"
          />
        )}
      </Card>
    </motion.div>
  );
}