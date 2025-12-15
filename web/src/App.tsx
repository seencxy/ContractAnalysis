import { ConfigProvider, theme } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { RouterProvider } from 'react-router-dom';
import { router } from './router';

function App() {
  return (
    <ConfigProvider
      locale={zhCN}
      theme={{
        token: {
          colorPrimary: '#1677ff',
          borderRadius: 12,
          fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
          colorBgLayout: '#f0f2f5',
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.05)',
        },
        components: {
          Card: {
            headerFontSize: 16,
            actionsBg: '#fff',
            boxShadowTertiary: '0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02)',
          },
          Table: {
            headerBg: 'transparent',
            headerColor: '#8c8c8c',
            headerSplitColor: 'transparent',
            rowHoverBg: '#fafafa',
          },
          Layout: {
            bodyBg: '#f0f2f5',
            headerBg: '#ffffff',
            siderBg: '#001529',
          },
          Menu: {
            itemBorderRadius: 8,
            itemMarginInline: 8,
          },
          Button: {
            borderRadius: 8,
            controlHeight: 36,
          },
        },
        algorithm: theme.defaultAlgorithm,
      }}
    >
      <RouterProvider router={router} />
    </ConfigProvider>
  );
}

export default App;
