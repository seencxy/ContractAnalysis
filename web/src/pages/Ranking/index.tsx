import { useState } from 'react';
import { Card, Radio, Table, Space, Typography, Tag, Button, Modal, Row, Col } from 'antd';
import type { RadioChangeEvent } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { TrophyOutlined, RiseOutlined, FallOutlined, FireOutlined } from '@ant-design/icons';
import { useSymbolStatistics } from '@/hooks/queries/useStatistics';
import type { Statistics } from '@/types/statistics';
import { Loading } from '@/components/common/Loading';
import { EmptyState } from '@/components/common/EmptyState';
import { motion } from 'framer-motion';

const { Title, Text } = Typography;

type PeriodType = '24h' | '7d' | '30d' | 'all';

export default function Ranking() {
  const [period, setPeriod] = useState<PeriodType>('24h');
  const [selectedSymbol, setSelectedSymbol] = useState<Statistics | null>(null);
  const [modalVisible, setModalVisible] = useState(false);

  const { data: response, isLoading } = useSymbolStatistics({ period });
  const symbols = response?.data || [];

  // 按胜率排序
  const sortedSymbols = [...symbols].sort((a, b) => {
    const aRate = a.win_rate ? parseFloat(a.win_rate) : 0;
    const bRate = b.win_rate ? parseFloat(b.win_rate) : 0;
    return bRate - aRate;
  });

  const handlePeriodChange = (e: RadioChangeEvent) => {
    setPeriod(e.target.value);
  };

  const showDetail = (record: Statistics) => {
    setSelectedSymbol(record);
    setModalVisible(true);
  };

  const columns: ColumnsType<Statistics> = [
    {
      title: '排名',
      key: 'rank',
      width: 80,
      align: 'center',
      render: (_: unknown, __: Statistics, index: number) => {
        let icon;
        if (index === 0) icon = <TrophyOutlined style={{ fontSize: 24, color: '#fadb14' }} />;
        else if (index === 1) icon = <TrophyOutlined style={{ fontSize: 22, color: '#d9d9d9' }} />;
        else if (index === 2) icon = <TrophyOutlined style={{ fontSize: 20, color: '#d48806' }} />;
        else icon = <div style={{ width: 24, height: 24, borderRadius: '50%', background: '#f5f5f5', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto', fontWeight: 600, color: '#8c8c8c' }}>{index + 1}</div>;
        
        return icon;
      },
    },
    {
      title: '交易对',
      dataIndex: 'symbol',
      key: 'symbol',
      width: 140,
      render: (symbol?: string) => (
        <Text strong style={{ fontSize: 15 }}>
          {symbol || '-'}
        </Text>
      ),
    },
    {
      title: '信号数',
      dataIndex: 'total_signals',
      key: 'total_signals',
      align: 'center',
      width: 100,
      sorter: (a, b) => a.total_signals - b.total_signals,
      render: (val) => <Tag bordered={false}>{val}</Tag>,
    },
    {
      title: '胜率',
      dataIndex: 'win_rate',
      key: 'win_rate',
      align: 'center',
      width: 120,
      render: (value?: string) => {
        if (!value) return '-';
        const rate = parseFloat(value);
        let color = 'default';
        if (rate >= 60) color = 'success';
        else if (rate >= 40) color = 'warning';
        else color = 'error';

        return (
          <Tag color={color} style={{ fontSize: 13, padding: '2px 8px', fontWeight: 600 }}>
            {rate.toFixed(1)}%
          </Tag>
        );
      },
      sorter: (a, b) => {
        const aRate = a.win_rate ? parseFloat(a.win_rate) : 0;
        const bRate = b.win_rate ? parseFloat(b.win_rate) : 0;
        return aRate - bRate;
      },
      defaultSortOrder: 'descend',
    },
    {
      title: '平均收益率',
      dataIndex: 'avg_profit_pct',
      key: 'avg_profit_pct',
      align: 'center',
      width: 150,
      render: (value?: string) => {
        if (!value) return '-';
        const profit = parseFloat(value);
        return (
          <Space>
            {profit > 0 ? (
              <RiseOutlined style={{ color: '#52c41a' }} />
            ) : (
              <FallOutlined style={{ color: '#ff4d4f' }} />
            )}
            <Text strong style={{ color: profit > 0 ? '#52c41a' : '#ff4d4f' }}>
              {profit > 0 ? '+' : ''}
              {profit.toFixed(2)}%
            </Text>
          </Space>
        );
      },
      sorter: (a, b) => {
        const aProfit = a.avg_profit_pct ? parseFloat(a.avg_profit_pct) : 0;
        const bProfit = b.avg_profit_pct ? parseFloat(b.avg_profit_pct) : 0;
        return aProfit - bProfit;
      },
    },
    {
      title: '极值表现',
      key: 'extremes',
      align: 'center',
      width: 180,
      render: (_: unknown, record: Statistics) => {
        const best = record.best_signal_pct ? parseFloat(record.best_signal_pct) : 0;
        const worst = record.worst_signal_pct ? parseFloat(record.worst_signal_pct) : 0;
        return (
            <div style={{ display: 'flex', flexDirection: 'column', fontSize: 12 }}>
                <span style={{ color: '#52c41a' }}>Max: +{best.toFixed(2)}%</span>
                <span style={{ color: '#ff4d4f' }}>Min: {worst.toFixed(2)}%</span>
            </div>
        )
      }
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      align: 'center',
      render: (_: unknown, record: Statistics) => (
        <Button type="link" size="small" onClick={() => showDetail(record)}>
          分析
        </Button>
      ),
    },
  ];

  const MetricItem = ({ label, value, valueColor, subValue }: { label: string, value: string | number, valueColor?: string, subValue?: string }) => (
    <div style={{ marginBottom: 16 }}>
      <Text type="secondary" style={{ fontSize: 12 }}>{label}</Text>
      <div style={{ display: 'flex', alignItems: 'baseline', gap: 8 }}>
        <Text strong style={{ fontSize: 18, color: valueColor }}>{value}</Text>
        {subValue && <Text type="secondary" style={{ fontSize: 12 }}>{subValue}</Text>}
      </div>
    </div>
  );

  return (
    <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4 }}
    >
      <Card bordered={false} style={{ marginBottom: 24, borderRadius: 12 }} bodyStyle={{ padding: '16px 24px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <FireOutlined style={{ color: '#ff4d4f', fontSize: 20 }} />
            <Title level={4} style={{ margin: 0 }}>交易对排行榜</Title>
          </div>
          <Radio.Group value={period} onChange={handlePeriodChange} buttonStyle="solid">
            <Radio.Button value="24h">24小时</Radio.Button>
            <Radio.Button value="7d">7天</Radio.Button>
            <Radio.Button value="30d">30天</Radio.Button>
            <Radio.Button value="all">全部</Radio.Button>
          </Radio.Group>
        </div>
      </Card>

      <Card bordered={false} style={{ borderRadius: 12 }} bodyStyle={{ padding: 0 }}>
        {isLoading ? (
            <Loading />
        ) : symbols.length === 0 ? (
            <EmptyState message="暂无交易对统计数据" />
        ) : (
            <Table
            columns={columns}
            dataSource={sortedSymbols}
            rowKey={(record) => `${record.symbol}-${record.period_label}`}
            pagination={{
                pageSize: 20,
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 个交易对`,
            }}
            size="middle"
            />
        )}
      </Card>

      <Modal
        title={
          <Space>
            <TrophyOutlined style={{ color: '#faad14' }} />
            <span style={{ fontSize: 18 }}>{selectedSymbol?.symbol} 详细数据</span>
          </Space>
        }
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={720}
        styles={{ body: { padding: '24px 24px 40px' } }}
      >
        {selectedSymbol && (
          <Space direction="vertical" size="large" style={{ width: '100%' }}>
            
            <div style={{ background: '#f5f7fa', padding: '20px', borderRadius: 8 }}>
              <Row gutter={[24, 24]}>
                <Col span={8}>
                  <MetricItem 
                    label="总信号数 / 已确认" 
                    value={`${selectedSymbol.total_signals} / ${selectedSymbol.confirmed_signals}`} 
                  />
                </Col>
                <Col span={8}>
                   <MetricItem 
                    label="胜率" 
                    value={selectedSymbol.win_rate ? `${parseFloat(selectedSymbol.win_rate).toFixed(2)}%` : '-'}
                    valueColor={selectedSymbol.win_rate && parseFloat(selectedSymbol.win_rate) >= 50 ? '#52c41a' : '#ff4d4f'}
                  />
                </Col>
                <Col span={8}>
                  <MetricItem 
                    label="盈亏比" 
                    value={selectedSymbol.profit_factor ? parseFloat(selectedSymbol.profit_factor).toFixed(2) : '-'}
                  />
                </Col>
                <Col span={8}>
                   <MetricItem 
                    label="盈利 / 亏损 信号" 
                    value={`${selectedSymbol.profitable_signals} / ${selectedSymbol.losing_signals}`}
                    valueColor="#666"
                  />
                </Col>
                <Col span={8}>
                   <MetricItem 
                    label="平均盈利" 
                    value={selectedSymbol.avg_profit_pct ? `+${parseFloat(selectedSymbol.avg_profit_pct).toFixed(2)}%` : '-'}
                    valueColor="#52c41a"
                  />
                </Col>
                <Col span={8}>
                   <MetricItem 
                    label="平均亏损" 
                    value={selectedSymbol.avg_loss_pct ? `${parseFloat(selectedSymbol.avg_loss_pct).toFixed(2)}%` : '-'}
                    valueColor="#ff4d4f"
                  />
                </Col>
              </Row>
            </div>

            <Row gutter={[32, 32]}>
              <Col span={12}>
                <Title level={5} style={{ marginBottom: 16 }}>K线分析</Title>
                <Row gutter={[16, 16]}>
                  <Col span={24}>
                    <MetricItem 
                        label="总K线小时数 / 最高价盈利小时" 
                        value={`${selectedSymbol.total_kline_hours} / ${selectedSymbol.profitable_kline_hours_high}`} 
                      />
                  </Col>
                  <Col span={24}>
                    <MetricItem 
                        label="K线最高价胜率" 
                        value={selectedSymbol.kline_theoretical_win_rate ? `${parseFloat(selectedSymbol.kline_theoretical_win_rate).toFixed(2)}%` : '-'} 
                        valueColor={selectedSymbol.kline_theoretical_win_rate && parseFloat(selectedSymbol.kline_theoretical_win_rate) >= 50 ? '#52c41a' : undefined}
                      />
                  </Col>
                  <Col span={24}>
                     <MetricItem 
                        label="K线收盘价胜率" 
                        value={selectedSymbol.kline_close_win_rate ? `${parseFloat(selectedSymbol.kline_close_win_rate).toFixed(2)}%` : '-'} 
                        valueColor={selectedSymbol.kline_close_win_rate && parseFloat(selectedSymbol.kline_close_win_rate) >= 50 ? '#52c41a' : undefined}
                      />
                  </Col>
                </Row>
              </Col>
              
              <Col span={12} style={{ borderLeft: '1px solid #f0f0f0', paddingLeft: 32 }}>
                <Title level={5} style={{ marginBottom: 16 }}>其他指标</Title>
                <Row gutter={[16, 16]}>
                   <Col span={24}>
                    <MetricItem 
                      label="平均持仓时长" 
                      value={selectedSymbol.avg_holding_hours ? `${parseFloat(selectedSymbol.avg_holding_hours).toFixed(1)} 小时` : '-'}
                    />
                   </Col>
                   <Col span={24}>
                    <div style={{ marginTop: 8 }}>
                      <Text type="secondary" style={{ fontSize: 12 }}>统计周期</Text>
                      <div style={{ marginTop: 4 }}>
                         <Tag>{selectedSymbol.period_label}</Tag>
                      </div>
                    </div>
                   </Col>
                </Row>
              </Col>
            </Row>

          </Space>
        )}
      </Modal>
    </motion.div>
  );
}