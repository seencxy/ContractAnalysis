# Binance Futures Analysis System

币安合约数据采集与分析系统 - 基于 Go 语言的智能交易信号生成系统

## 🎯 功能特性

- **自动数据采集**：定时采集所有 USDT 本位合约的多空比数据
- **智能策略分析**：
  - 逆向策略（Minority Strategy）：跟随少数派方向
  - 大户策略（Whale Strategy）：分析持仓与账户数量的分离
- **信号追踪**：监控信号发出后的价格走势，统计盈利概率
- **灵活配置**：所有策略参数支持配置文件修改
- **多通道通知**：支持控制台、Telegram、邮件等通知方式

## 📋 系统要求

- Go 1.25+
- MySQL 5.7+ 或 MariaDB 10.3+
- Redis 6.0+（可选，用于缓存）

## 🚀 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置数据库

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE futures_analysis CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 运行数据库迁移
mysql -u root -p futures_analysis < scripts/migrations/001_initial_schema.sql
```

### 3. 配置系统

编辑 `config.yaml` 文件：

```yaml
# 数据库配置
database:
  mysql:
    host: "localhost"
    port: 3306
    database: "futures_analysis"
    username: "root"
    password: "your_password"  # 修改为你的密码
  redis:
    host: "localhost"
    port: 6379

# Binance API 配置（如果需要实时数据）
binance:
  api_key: ""     # 可选，某些端点不需要
  api_secret: ""  # 可选，某些端点不需要
```

### 4. 运行系统

```bash
go run main.go
```

或编译后运行：

```bash
go build -o futures-analysis
./futures-analysis
```

## 🐳 Docker 一键部署

### 快速启动

```bash
# 1. 初始化环境配置
./scripts/docker/deploy.sh init

# 2. 编辑配置文件（修改密码等）
vim .env.docker

# 3. 一键启动所有服务
./scripts/docker/deploy.sh up
```

### 部署命令

```bash
# 启动服务
./scripts/docker/deploy.sh up

# 停止服务
./scripts/docker/deploy.sh down

# 重启服务
./scripts/docker/deploy.sh restart

# 查看日志
./scripts/docker/deploy.sh logs -f

# 查看服务状态
./scripts/docker/deploy.sh status

# 重新构建
./scripts/docker/deploy.sh build

# 清理所有数据（危险操作）
./scripts/docker/deploy.sh clean
```

### Docker Compose 直接使用

```bash
# 创建环境配置
cp .env.docker.example .env.docker

# 启动
docker compose --env-file .env.docker up -d

# 查看日志
docker compose logs -f

# 停止
docker compose down
```

### 服务说明

| 服务 | 端口 | 说明 |
|------|------|------|
| app | 8080 | 主应用 API |
| metrics | 9090 | Prometheus 指标 |
| mysql | 3306 | MySQL 数据库 |
| redis | 6379 | Redis 缓存 |

### 数据持久化

- MySQL 数据：`mysql_data` Docker Volume
- Redis 数据：`redis_data` Docker Volume
- 应用日志：`./logs/` 目录

## 📊 系统架构

```
ContractAnalysis/
├── config/                      # 配置系统
├── internal/
│   ├── domain/                  # 领域层
│   │   ├── entity/              # 领域实体
│   │   ├── repository/          # 仓储接口
│   │   └── service/             # 策略服务
│   ├── usecase/                 # 业务用例
│   │   ├── collector.go         # 数据采集
│   │   ├── analyzer.go          # 信号分析
│   │   └── tracker.go           # 信号追踪
│   └── infrastructure/          # 基础设施
│       ├── binance/             # Binance API
│       ├── persistence/         # 数据持久化
│       ├── notification/        # 通知系统
│       ├── scheduler/           # 定时调度
│       └── logger/              # 日志系统
└── main.go                      # 应用入口
```

## 🎮 使用说明

### 定时任务

系统默认配置了以下定时任务：

- **数据采集**：每小时执行一次（可在 `config.yaml` 中配置）
- **信号分析**：每小时第5分钟执行
- **信号追踪**：每15分钟执行一次

### 策略配置

#### 逆向策略（Minority Strategy）

跟随少数派方向：当 80% 账户做空时，生成做多信号。

```yaml
strategies:
  minority:
    enabled: true
    min_ratio_difference: 75.0        # 多空比差距 >= 75:25
    confirmation_hours: 2             # 需要2小时确认
    generate_long_when_short_ratio_above: 75.0
    generate_short_when_long_ratio_above: 75.0
    tracking_hours: 24                # 追踪24小时
    profit_target_pct: 5.0            # 目标盈利5%
    stop_loss_pct: 2.0                # 止损2%
```

#### 大户策略（Whale Strategy）

分析持仓与账户数量的分离度，识别散户被收割场景。

```yaml
strategies:
  whale:
    enabled: true
    min_ratio_difference: 75.0        # 账户比例差距
    whale_position_threshold: 60.0    # 大户持仓占比 >= 60%
    min_divergence: 20.0              # 账户比与持仓比最小分离度
    confirmation_hours: 2
    tracking_hours: 24
```

### 通知配置

#### 控制台通知

```yaml
notifications:
  console:
    enabled: true
    events:
      - "signal_generated"
      - "signal_confirmed"
      - "signal_outcome"
```

#### Telegram 通知（可选）

```yaml
notifications:
  telegram:
    enabled: true
    bot_token: "YOUR_BOT_TOKEN"
    chat_ids:
      - "YOUR_CHAT_ID"
    events:
      - "signal_generated"
      - "signal_confirmed"
```

## 📈 信号生命周期

1. **生成（PENDING）**：策略检测到符合条件的市场状态
2. **确认（CONFIRMED）**：经过确认期（默认2小时）后，条件仍然满足
3. **追踪（TRACKING）**：开始追踪价格变化
4. **关闭（CLOSED）**：达到止盈/止损或追踪期结束

## 🔍 数据查询

### 查看最新信号

```sql
SELECT signal_id, symbol, signal_type, strategy_name, status,
       price_at_signal, generated_at, reason
FROM signals
ORDER BY generated_at DESC
LIMIT 10;
```

### 查看信号结果统计

```sql
SELECT strategy_name, outcome, COUNT(*) as count,
       AVG(final_price_change_pct) as avg_change
FROM signal_outcomes so
JOIN signals s ON so.signal_id = s.signal_id
GROUP BY strategy_name, outcome;
```

### 查看策略胜率

```sql
SELECT strategy_name, period_label,
       win_rate, avg_profit_pct, avg_loss_pct,
       total_signals, profitable_signals, losing_signals
FROM strategy_statistics
WHERE period_label = '24h'
ORDER BY calculated_at DESC;
```

## ⚙️ 环境变量

可以通过环境变量覆盖配置文件中的设置：

```bash
# 数据库配置
export CA_DATABASE_MYSQL_HOST=localhost
export CA_DATABASE_MYSQL_PORT=3306
export CA_DATABASE_MYSQL_USERNAME=root
export CA_DATABASE_MYSQL_PASSWORD=your_password

# Binance API（可选）
export CA_BINANCE_API_KEY=your_api_key
export CA_BINANCE_API_SECRET=your_api_secret

# Redis 配置
export CA_DATABASE_REDIS_HOST=localhost
export CA_DATABASE_REDIS_PORT=6379
```

## 🛠️ 开发指南

### 添加新策略

1. 在 `internal/domain/service/` 创建新策略文件
2. 实现 `Strategy` 接口
3. 在 `main.go` 中注册新策略

示例：

```go
type MyCustomStrategy struct {
    *BaseStrategy
    config MyCustomStrategyConfig
}

func (s *MyCustomStrategy) Analyze(ctx context.Context, recentData []*entity.MarketData) ([]*entity.Signal, error) {
    // 实现你的策略逻辑
    return signals, nil
}
```

### 添加新通知渠道

1. 在 `internal/infrastructure/notification/` 创建新通知器
2. 实现 `Notifier` 接口
3. 在 `main.go` 中注册新通知器

## 📝 日志

日志文件位置：`logs/app.log`

日志级别：
- `debug`：详细调试信息
- `info`：一般信息
- `warn`：警告信息
- `error`：错误信息

修改日志级别：

```yaml
logging:
  level: "info"  # debug, info, warn, error
  format: "json" # json or console
```

## 🔒 安全建议

1. **保护 API 密钥**：不要将 Binance API 密钥提交到代码仓库
2. **数据库安全**：使用强密码，限制数据库访问权限
3. **仅分析模式**：系统默认仅进行分析和提醒，不执行实际交易
4. **备份数据**：定期备份数据库数据

## 🐛 故障排查

### 数据库连接失败

检查：
1. MySQL 服务是否启动
2. 数据库配置是否正确
3. 用户权限是否足够

### Binance API 失败

检查：
1. 网络连接是否正常
2. API 密钥是否正确（如果使用）
3. 是否触发了 API 限流

### 信号不生成

检查：
1. 策略是否启用（`enabled: true`）
2. 市场条件是否满足策略阈值
3. 是否在冷却期内

## 📊 性能优化

### 数据库优化

```sql
-- 定期清理旧数据（保留30天）
DELETE FROM market_data WHERE timestamp < DATE_SUB(NOW(), INTERVAL 30 DAY);

-- 优化表
OPTIMIZE TABLE market_data;
OPTIMIZE TABLE signals;
OPTIMIZE TABLE signal_tracking;
```

### Redis 缓存

启用 Redis 可以提高查询性能：

```yaml
database:
  type: "mysql"
  redis:
    enabled: true
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 📞 支持

如有问题，请提交 Issue 或联系开发团队。

---

**免责声明**：本系统仅供学习和研究使用，不构成投资建议。使用本系统进行交易的风险由用户自行承担。
