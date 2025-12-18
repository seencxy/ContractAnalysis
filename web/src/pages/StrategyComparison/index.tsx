import { useState } from 'react';
import { Select, Card, Table, Tag, Row, Col, Typography, Space } from 'antd';
import { useStrategyComparison } from '@/hooks/queries/useStatistics';
import { WinRateChart } from '@/components/charts/WinRateChart';
import { Loading } from '@/components/common/Loading';
import { EmptyState } from '@/components/common/EmptyState';
import type { Statistics } from '@/types/statistics';

const { Title, Text } = Typography;
const { Option } = Select;

function getStrategyDisplayName(key: string): string {
  const names: Record<string, string> = {
    'MinorityStrategy': '少数派策略',
    'WhaleStrategy': '鲸鱼策略',
    'SmartMoneyStrategy': '聪明钱策略'
  };
  return names[key] || key;
}

export default function StrategyComparison() {
  const [period, setPeriod] = useState<'24h' | '7d' | '30d' | 'all'>('24h');
  const [selectedStrategies, setSelectedStrategies] = useState<string[]>([
    'MinorityStrategy',
    'WhaleStrategy',
    'SmartMoneyStrategy'
  ]);

  const { data: response, isLoading, error } = useStrategyComparison({
    strategies: selectedStrategies,
    period
  });

  const comparisonData = response?.data;

  // Table columns
  const columns = [
    {
      title: '策略',
      dataIndex: 'strategy_name',
      key: 'strategy_name',
      render: (name: string) => (
        <Space>
          <Text strong>{getStrategyDisplayName(name)}</Text>
          {comparisonData?.comparison?.best_win_rate === name && (
            <Tag color="gold">最佳胜率</Tag>
          )}
          {comparisonData?.comparison?.best_avg_return === name && (
            <Tag color="green">最佳收益</Tag>
          )}
        </Space>
      )
    },
    {
      title: '信号数',
      dataIndex: 'total_signals',
      key: 'total_signals',
      sorter: (a: Statistics, b: Statistics) => a.total_signals - b.total_signals
    },
    {
      title: '胜率',
      dataIndex: 'win_rate',
      key: 'win_rate',
      render: (rate: string) => `${rate}%`,
      sorter: (a: Statistics, b: Statistics) => parseFloat(a.win_rate || '0') - parseFloat(b.win_rate || '0')
    },
    {
      title: '平均收益',
      dataIndex: 'avg_profit_pct',
      key: 'avg_profit_pct',
      render: (pct: string) => (
        <Text type={parseFloat(pct) > 0 ? 'success' : 'danger'}>
          {pct}%
        </Text>
      ),
      sorter: (a: Statistics, b: Statistics) => parseFloat(a.avg_profit_pct || '0') - parseFloat(b.avg_profit_pct || '0')
    },
    {
      title: '盈亏比',
      dataIndex: 'profit_factor',
      key: 'profit_factor',
      sorter: (a: Statistics, b: Statistics) => parseFloat(a.profit_factor || '0') - parseFloat(b.profit_factor || '0')
    },
    {
      title: '盈利 / 亏损',
      key: 'counts',
      render: (_: any, record: Statistics) => (
        <Space>
          <Text type="success">{record.profitable_signals}</Text>
          <Text type="secondary">/</Text>
          <Text type="danger">{record.losing_signals}</Text>
        </Space>
      )
    }
  ];

  return (
    <div style={{ paddingBottom: 24 }}>
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Title level={4} style={{ margin: 0 }}>策略对比分析</Title>
        <Space>
          <Select
            mode="multiple"
            value={selectedStrategies}
            onChange={setSelectedStrategies}
            style={{ width: 400 }}
            placeholder="选择策略（2-5个）"
            maxTagCount="responsive"
          >
            <Option value="MinorityStrategy">少数派策略</Option>
            <Option value="WhaleStrategy">鲸鱼策略</Option>
            <Option value="SmartMoneyStrategy">聪明钱策略</Option>
          </Select>
          <Select
            value={period}
            onChange={setPeriod}
            style={{ width: 120 }}
          >
            <Option value="24h">24小时</Option>
            <Option value="7d">7天</Option>
            <Option value="30d">30天</Option>
            <Option value="all">全部</Option>
          </Select>
        </Space>
      </div>

      {isLoading ? (
        <Loading />
      ) : error ? (
        <EmptyState message="加载失败，请稍后重试" />
      ) : !comparisonData ? (
        <EmptyState message="暂无对比数据" />
      ) : (
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          {/* Summary Cards */}
          <Row gutter={16}>
            <Col span={8}>
              <Card size="small" bordered={false} style={{ background: '#f6ffed', borderRadius: 8 }}>
                <Text type="secondary">最佳胜率</Text>
                <div style={{ fontSize: 20, fontWeight: 600, marginTop: 4 }}>
                  {comparisonData.comparison.best_win_rate ? getStrategyDisplayName(comparisonData.comparison.best_win_rate) : '--'}
                  <Text type="success" style={{ marginLeft: 8 }}>
                    {comparisonData.comparison.best_win_rate ? comparisonData.comparison.win_rates[comparisonData.comparison.best_win_rate] : '--'}%
                  </Text>
                </div>
              </Card>
            </Col>
            <Col span={8}>
              <Card size="small" bordered={false} style={{ background: '#f0f5ff', borderRadius: 8 }}>
                <Text type="secondary">最佳收益</Text>
                <div style={{ fontSize: 20, fontWeight: 600, marginTop: 4 }}>
                  {comparisonData.comparison.best_avg_return ? getStrategyDisplayName(comparisonData.comparison.best_avg_return) : '--'}
                  <Text type="success" style={{ marginLeft: 8 }}>
                    {comparisonData.comparison.best_avg_return ? comparisonData.comparison.avg_returns[comparisonData.comparison.best_avg_return] : '--'}%
                  </Text>
                </div>
              </Card>
            </Col>
            <Col span={8}>
              <Card size="small" bordered={false} style={{ background: '#fff7e6', borderRadius: 8 }}>
                <Text type="secondary">信号最多</Text>
                <div style={{ fontSize: 20, fontWeight: 600, marginTop: 4 }}>
                  {comparisonData.comparison.most_signals ? getStrategyDisplayName(comparisonData.comparison.most_signals) : '--'}
                  <Text type="warning" style={{ marginLeft: 8 }}>
                    {comparisonData.comparison.most_signals && comparisonData.comparison.total_signals[comparisonData.comparison.most_signals] !== 0
                      ? comparisonData.comparison.total_signals[comparisonData.comparison.most_signals]
                      : '--'}
                  </Text>
                </div>
              </Card>
            </Col>
          </Row>

          {/* Charts */}
          <Row gutter={24} style={{ height: 400 }}>
             <Col span={24} style={{ height: '100%' }}>
                <Card bordered={false} style={{ borderRadius: 12, height: '100%' }}>
                    {comparisonData.detailed_stats && comparisonData.detailed_stats.length > 0 ? (
                        <WinRateChart data={comparisonData.detailed_stats} />
                    ) : (
                        <EmptyState message="暂无图表数据" />
                    )}
                </Card>
             </Col>
          </Row>

          {/* Detailed Stats Table */}
          <Card bordered={false} style={{ borderRadius: 12 }}>
            <Table
              dataSource={comparisonData.detailed_stats}
              columns={columns}
              rowKey="strategy_name"
              pagination={false}
              locale={{ emptyText: <EmptyState message="暂无详细统计数据" /> }}
            />
          </Card>
        </Space>
      )}
    </div>
  );
}
