import { Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

export function Loading() {
  return (
    <div 
      style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        minHeight: 200, 
        width: '100%',
        padding: '24px'
      }}
    >
      <Spin 
        indicator={<LoadingOutlined style={{ fontSize: 32, color: '#1677ff' }} spin />} 
        tip={<span style={{ marginTop: 12, color: '#8c8c8c' }}>Loading...</span>}
      />
    </div>
  );
}