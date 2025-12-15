# ContractAnalysis React 前端

这是一个为 Binance Futures Analysis 系统构建的现代化 React 前端应用。

## 技术栈

- **框架**: React 18 + TypeScript + Vite
- **状态管理**: TanStack Query (React Query) + Zustand
- **UI 组件**: Ant Design 5
- **图表**: Apache ECharts
- **路由**: React Router v6
- **HTTP 客户端**: Axios

## 项目结构

```
web/
├── src/
│   ├── api/                          # API 调用层
│   │   ├── client.ts                # Axios 实例配置
│   │   └── endpoints/               # API 端点封装
│   │       ├── signals.ts           # 信号相关API
│   │       └── statistics.ts        # 统计相关API
│   ├── components/                  # 组件库
│   │   ├── layout/                 # 布局组件
│   │   │   └── MainLayout.tsx      # 主布局(侧边栏导航)
│   │   ├── common/                 # 通用组件
│   │   │   ├── Loading.tsx         # 加载组件
│   │   │   └── EmptyState.tsx      # 空状态组件
│   │   └── charts/                 # 图表组件
│   │       ├── WinRateChart.tsx            # 胜率对比图
│   │       ├── SignalStatusPieChart.tsx    # 状态分布图
│   │       └── PriceLineChart.tsx          # 价格走势图
│   ├── hooks/                      # 自定义 Hooks
│   │   └── queries/                # React Query hooks
│   │       ├── useSignals.ts       # 信号查询hooks
│   │       └── useStatistics.ts    # 统计查询hooks
│   ├── pages/                      # 页面组件
│   │   ├── Dashboard/              # 仪表盘
│   │   ├── Signals/                # 实时信号监控
│   │   ├── Analysis/               # 策略分析
│   │   ├── Ranking/                # 交易对排名
│   │   └── History/                # 历史查询
│   ├── types/                      # TypeScript 类型定义
│   │   ├── common.ts               # 通用类型
│   │   ├── signal.ts               # 信号类型
│   │   └── statistics.ts           # 统计类型
│   ├── utils/                      # 工具函数
│   │   ├── format.ts               # 格式化工具
│   │   └── colors.ts               # 颜色工具
│   ├── App.tsx                     # 根组件
│   ├── main.tsx                    # 入口文件
│   └── router.tsx                  # 路由配置
└── vite.config.ts                  # Vite 配置
```

## 快速开始

### 1. 安装依赖

```bash
cd web
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

前端将在 http://localhost:5173 启动

### 3. 构建生产版本

```bash
npm run build
```

构建产物将输出到 `dist/` 目录

## 重要说明

### API 响应结构

后端所有 API 返回都使用统一的响应包装格式:

```typescript
{
  code: 200,
  message: "success",
  data: { ... },  // 实际数据
  timestamp: 1234567890
}
```

在 React Query hooks 中使用时需要访问 `.data` 属性:

```typescript
const { data: response } = useSignals({ page: 1, limit: 20 });
const signals = response?.data?.items || [];
```

### 配置

- **Vite 代理**: 自动将 `/api` 请求代理到 `http://localhost:8080`
- **TypeScript 路径别名**: 使用 `@/` 引用 `src/` 目录

## 已实现功能

### ✅ 仪表盘页面 (`/`)
- 关键指标卡片展示(今日信号、活跃信号、胜率、收益)
- 信号状态分布饼图
- 最近信号列表

### ✅ 实时信号监控页面 (`/signals`)
- 信号列表展示
- 状态、类型、交易对筛选
- 分页功能
- 实时刷新
- 信号详情抽屉(包含价格走势图和追踪数据)

### ✅ 策略分析页面 (`/analysis`)
- 时间周期选择器(24h/7d/30d/all)
- 策略胜率对比图表
- 策略详细数据表格
- K线分析指标
- 理论收益分析

### ✅ 交易对排名页面 (`/ranking`)
- 交易对表现排名表格
- 多维度排序(胜率、信号数、平均收益)
- 交易对详情模态框
- K线胜率数据展示

### ✅ 历史查询页面 (`/history`)
- 高级筛选面板(时间范围、交易对、类型、状态、策略)
- 信号历史记录表格
- 分页和快速跳转
- 信号详情抽屉(包含追踪记录时间轴)

### ✅ 图表组件
- 胜率对比柱状图 (WinRateChart)
- 信号状态分布饼图 (SignalStatusPieChart)
- 价格走势折线图 (PriceLineChart)

## 开发建议

1. **代码分割**: 对于大型页面组件，使用 React.lazy() 进行代码分割
2. **性能优化**: 合理设置 React Query 的 staleTime 和 refetchInterval
3. **错误处理**: 完善全局错误边界和 API 错误处理
4. **响应式设计**: 添加移动端适配

## API 端点

前端已封装以下 API 端点:

**信号相关:**
- `GET /api/v1/signals` - 信号列表
- `GET /api/v1/signals/active` - 活跃信号
- `GET /api/v1/signals/:id` - 信号详情
- `GET /api/v1/signals/:id/tracking` - 追踪记录
- `GET /api/v1/signals/:id/klines` - K线数据

**统计相关:**
- `GET /api/v1/statistics/overview` - 统计概览
- `GET /api/v1/statistics/strategies` - 策略统计
- `GET /api/v1/statistics/symbols` - 交易对统计

## 许可证

MIT
