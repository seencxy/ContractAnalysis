import { Layout, Menu, Button, Avatar, Dropdown, Space, theme } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  DashboardOutlined,
  SignalFilled,
  LineChartOutlined,
  TrophyOutlined,
  HistoryOutlined,
  BarChartOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  BellOutlined,
} from '@ant-design/icons';
import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

const { Header, Sider, Content } = Layout;

export function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const [collapsed, setCollapsed] = useState(false);
  const { token } = theme.useToken();

  const menuItems = [
    { key: '/', label: '仪表盘', icon: <DashboardOutlined /> },
    { key: '/signals', label: '实时信号', icon: <SignalFilled /> },
    { key: '/analysis', label: '策略分析', icon: <LineChartOutlined /> },
    { key: '/ranking', label: '交易对排名', icon: <TrophyOutlined /> },
    { key: '/history', label: '历史查询', icon: <HistoryOutlined /> },
    { key: '/statistics-history', label: '统计历史', icon: <BarChartOutlined /> },
  ];

  return (
    <Layout style={{ minHeight: '100vh', background: '#f0f2f5' }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        width={240}
        theme="light"
        style={{
          boxShadow: '2px 0 8px 0 rgba(29,35,41,.05)',
          zIndex: 10,
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
          height: '100vh',
          overflow: 'auto',
        }}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            borderBottom: '1px solid rgba(5, 5, 5, 0.06)',
            marginBottom: 16,
          }}
        >
          <motion.div
            initial={false}
            animate={{ scale: collapsed ? 0.8 : 1 }}
            style={{ 
              color: token.colorPrimary, 
              fontSize: collapsed ? 18 : 20, 
              fontWeight: 800,
              display: 'flex',
              alignItems: 'center',
              gap: 8,
              whiteSpace: 'nowrap',
              overflow: 'hidden'
            }}
          >
            <SignalFilled />
            {!collapsed && "Futures Analysis"}
          </motion.div>
        </div>
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          style={{ borderRight: 0 }}
        />
      </Sider>
      
      <Layout style={{ marginLeft: collapsed ? 80 : 240, transition: 'all 0.2s' }}>
        <Header
          className="glass-effect"
          style={{
            padding: '0 24px',
            position: 'sticky',
            top: 0,
            zIndex: 9,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            height: 64,
          }}
        >
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{ fontSize: '16px', width: 46, height: 46 }}
          />
          
          <Space size={16}>
            <Button type="text" icon={<BellOutlined />} style={{ borderRadius: '50%', width: 40, height: 40 }} />
            <Dropdown menu={{ items: [{ key: '1', label: '个人设置' }, { key: '2', label: '退出登录' }] }}>
              <Space style={{ cursor: 'pointer' }}>
                <Avatar icon={<UserOutlined />} style={{ backgroundColor: token.colorPrimary }} />
                <span style={{ fontWeight: 500 }}>Admin</span>
              </Space>
            </Dropdown>
          </Space>
        </Header>
        
        <Content style={{ margin: '24px 24px 0', minHeight: 280, overflow: 'hidden' }}>
          <AnimatePresence mode="wait">
            <motion.div
              key={location.pathname}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3, ease: "easeInOut" }}
            >
              <Outlet />
            </motion.div>
          </AnimatePresence>
        </Content>
        <Layout.Footer style={{ textAlign: 'center', color: '#8c8c8c' }}>
          Binance Futures Analysis System ©{new Date().getFullYear()}
        </Layout.Footer>
      </Layout>
    </Layout>
  );
}