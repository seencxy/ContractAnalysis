import { Empty } from 'antd';

interface EmptyStateProps {
  message?: string;
}

export function EmptyState({ message = '暂无数据' }: EmptyStateProps) {
  return (
    <div style={{ padding: '48px 0', display: 'flex', justifyContent: 'center' }}>
      <Empty 
        image={Empty.PRESENTED_IMAGE_SIMPLE} 
        description={<span style={{ color: '#8c8c8c' }}>{message}</span>} 
      />
    </div>
  );
}