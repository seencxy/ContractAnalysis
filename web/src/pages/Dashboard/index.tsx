import { Row, Col, Card, Table, Tag, Typography } from 'antd';
import {
  SignalFilled,
  ThunderboltFilled,
  RiseOutlined,
  TrophyOutlined,
  ArrowRightOutlined,
} from '@ant-design/icons';
import { useOverviewStatistics } from '@/hooks/queries/useStatistics';
import { useSignals } from '@/hooks/queries/useSignals';
import { Loading } from '@/components/common/Loading';
import { EmptyState } from '@/components/common/EmptyState';
import { SignalStatusPieChart } from '@/components/charts/SignalStatusPieChart';
import StrategyBadge from '@/components/common/StrategyBadge';
import { useStrategies } from '@/hooks/queries/useStrategies';
import { formatPercentString, formatRelativeTime } from '@/utils/format';
import { getStatusColor, getSignalTypeColor } from '@/utils/colors';
import { motion } from 'framer-motion';
import { useNavigate } from 'react-router-dom';

const { Title } = Typography;

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1
    }
  }
};

const itemVariants: any = {
  hidden: { y: 20, opacity: 0 },
  visible: {
    y: 0,
    opacity: 1,
    transition: {
      type: "spring",
      stiffness: 100
    }
  }
};

interface StatCardProps {
  title: string;
  value: string | number;
  prefix: React.ReactNode;
  color: string;
  loading?: boolean;
}

const StatCard = ({ title, value, prefix, color, loading }: StatCardProps) => (
  <Card 
    bordered={false} 
    className="card-hover-effect"
    style={{ height: '100%', borderRadius: 12, overflow: 'hidden' }}
    bodyStyle={{ padding: '20px 24px' }}
  >
    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
      <div>
        <div style={{ color: '#8c8c8c', fontSize: 14, marginBottom: 4 }}>{title}</div>
        <div style={{ fontSize: 28, fontWeight: 700, color: '#1f1f1f' }}>
          {loading ? '-' : value}
        </div>
      </div>
      <div 
        style={{ 
          width: 48, 
          height: 48, 
          borderRadius: 12, 
          background: `${color}15`, 
          display: 'flex', 
          alignItems: 'center', 
          justifyContent: 'center',
          color: color,
          fontSize: 24
        }}
      >
        {prefix}
      </div>
    </div>
  </Card>
);

export default function Dashboard() {
  const navigate = useNavigate();
  const { data: overviewRes, isLoading: overviewLoading } = useOverviewStatistics();
  const { data: recentRes, isLoading: recentLoading } = useSignals({ limit: 10, page: 1 });
  const { data: strategies, isLoading: strategiesLoading } = useStrategies();

  const overview = overviewRes?.data;
  const recentSignals = recentRes?.data?.items || [];

  console.log('Dashboard strategies:', strategies);
  console.log('Dashboard strategiesLoading:', strategiesLoading);

  if (overviewLoading) return <Loading />;

  const columns = [
    {
      title: '交易对',
      dataIndex: 'symbol',
      key: 'symbol',
      render: (text: string) => <span style={{ fontWeight: 600, color: '#1f1f1f' }}>{text}</span>,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color={getSignalTypeColor(type)} style={{ borderRadius: 6, padding: '0 8px' }}>{type}</Tag>
      ),
    },
    {
      title: '策略',
      dataIndex: 'strategy_name',
      key: 'strategy_name',
      render: (_: string, record: any) => <StrategyBadge signal={record} />,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={getStatusColor(status)} style={{ borderRadius: 6 }}>{status}</Tag>
      ),
    },
    {
      title: '生成时间',
      dataIndex: 'generated_at',
      key: 'generated_at',
      render: (time: string) => <span style={{ color: '#999' }}>{formatRelativeTime(time)}</span>,
    },
  ];

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
    >
      <div style={{ marginBottom: 24 }}>
        <Title level={4} style={{ margin: 0 }}>仪表盘概览</Title>
      </div>

      <Row gutter={[24, 24]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} md={6}>
          <motion.div variants={itemVariants}>
            <StatCard
              title="今日信号数"
              value={overview?.total_signals_today || 0}
              prefix={<SignalFilled />}
              color="#1677ff"
              loading={overviewLoading}
            />
          </motion.div>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <motion.div variants={itemVariants}>
            <StatCard
              title="活跃信号"
              value={overview?.active_signals || 0}
              prefix={<ThunderboltFilled />}
              color="#52c41a"
              loading={overviewLoading}
            />
          </motion.div>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <motion.div variants={itemVariants}>
            <StatCard
              title="24h胜率"
              value={formatPercentString(overview?.overall_win_rate_24h)}
              prefix={<RiseOutlined />}
              color="#faad14"
              loading={overviewLoading}
            />
          </motion.div>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <motion.div variants={itemVariants}>
            <StatCard
              title="24h平均收益"
              value={formatPercentString(overview?.avg_return_pct_24h)}
              prefix={<TrophyOutlined />}
              color="#722ed1"
              loading={overviewLoading}
            />
          </motion.div>
        </Col>
      </Row>

      <Row gutter={[24, 24]} style={{ marginBottom: 24 }}>
        <Col span={24}>
          <motion.div variants={itemVariants}>
            <Card title="策略实时表现 (24h)" bordered={false} style={{ borderRadius: 12 }}>
              {overviewLoading ? (
                <Loading />
              ) : (
                <Row gutter={[16, 16]}>
                  {overview?.strategy_breakdown?.map((s) => (
                    <Col xs={24} sm={8} key={s.strategy_name}>
                      <Card
                        size="small"
                        bordered={false}
                        style={{
                          borderRadius: 8,
                          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.05)',
                          transition: 'all 0.3s',
                          cursor: 'default',
                        }}
                      >
                        <div style={{ display: 'flex', flexDirection: 'column' }}>
                          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
                            <Typography.Text strong style={{ fontSize: 16 }}>
                              {s.strategy_name}
                            </Typography.Text>
                            <Tag color="blue">{s.signal_count} 信号</Tag>
                          </div>
                          <Row gutter={8}>
                            <Col span={12}>
                              <div style={{ color: '#8c8c8c', fontSize: 12 }}>胜率</div>
                              <div style={{ fontSize: 18, fontWeight: 600, color: '#1f1f1f' }}>
                                {s.win_rate || '-'}%
                              </div>
                            </Col>
                            <Col span={12}>
                              <div style={{ color: '#8c8c8c', fontSize: 12 }}>平均收益</div>
                              <div style={{ 
                                fontSize: 18, 
                                fontWeight: 600, 
                                color: parseFloat(s.avg_return_pct || '0') >= 0 ? '#52c41a' : '#ff4d4f' 
                              }}>
                                {s.avg_return_pct || '-'}%
                              </div>
                            </Col>
                          </Row>
                          <div style={{ marginTop: 12, display: 'flex', gap: 16, fontSize: 12 }}>
                            <span style={{ color: '#52c41a' }}>✅ {s.profitable_count} 盈利</span>
                            <span style={{ color: '#ff4d4f' }}>❌ {s.losing_count} 亏损</span>
                          </div>
                        </div>
                      </Card>
                    </Col>
                  ))}
                  {(!overview?.strategy_breakdown || overview.strategy_breakdown.length === 0) && (
                    <Col span={24}><EmptyState message="暂无策略表现数据" /></Col>
                  )}
                </Row>
              )}
            </Card>
          </motion.div>
        </Col>
      </Row>

      <Row gutter={[24, 24]}>
        <Col xs={24} lg={8}>
          <motion.div variants={itemVariants} style={{ height: '100%' }}>
            <Card 
              title="信号状态分布" 
              bordered={false}
              className="card-hover-effect"
              style={{ height: '100%', borderRadius: 12 }}
            >
              {overviewLoading ? (
                <Loading />
              ) : !overview?.status_distribution ? (
                <EmptyState message="暂无信号数据" />
              ) : (
                <div style={{ height: 300 }}>
                  <SignalStatusPieChart distribution={overview.status_distribution as unknown as Record<string, number>} />
                </div>
              )}
            </Card>
          </motion.div>
        </Col>

        <Col xs={24} lg={16}>
          <motion.div variants={itemVariants} style={{ height: '100%' }}>
            <Card 
              title="最近信号" 
              bordered={false}
              className="card-hover-effect"
              style={{ height: '100%', borderRadius: 12 }}
              extra={
                <div 
                  style={{ cursor: 'pointer', color: '#1677ff', display: 'flex', alignItems: 'center', gap: 4 }}
                  onClick={() => navigate('/signals')}
                >
                  查看全部 <ArrowRightOutlined />
                </div>
              }
            >
              {recentLoading ? (
                <Loading />
              ) : recentSignals.length === 0 ? (
                <EmptyState message="暂无信号数据" />
              ) : (
                <Table
                  dataSource={recentSignals}
                  columns={columns}
                  rowKey="signal_id"
                  pagination={false}
                  size="middle"
                />
              )}
            </Card>
          </motion.div>
        </Col>
      </Row>
    </motion.div>
  );
}