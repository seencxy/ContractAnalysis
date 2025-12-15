package config

import "time"

// Config represents the application configuration
type Config struct {
	App           AppConfig           `mapstructure:"app"`
	Server        ServerConfig        `mapstructure:"server"`
	Binance       BinanceConfig       `mapstructure:"binance"`
	Collection    CollectionConfig    `mapstructure:"collection"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Strategies    StrategiesConfig    `mapstructure:"strategies"`
	Statistics    StatisticsConfig    `mapstructure:"statistics"`
	Notifications NotificationsConfig `mapstructure:"notifications"`
	Logging       LoggingConfig       `mapstructure:"logging"`
	Monitoring    MonitoringConfig    `mapstructure:"monitoring"`
	Features      FeaturesConfig      `mapstructure:"features"`
}

// AppConfig represents general application configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Timezone    string `mapstructure:"timezone"`
}

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// BinanceConfig represents Binance API configuration
type BinanceConfig struct {
	APIURL    string          `mapstructure:"api_url"`
	APIKey    string          `mapstructure:"api_key"`
	APISecret string          `mapstructure:"api_secret"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Timeout   time.Duration   `mapstructure:"timeout"`
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"`
	WeightPerMinute   int `mapstructure:"weight_per_minute"`
}

// CollectionConfig represents data collection configuration
type CollectionConfig struct {
	Enabled    bool        `mapstructure:"enabled"`
	Interval   string      `mapstructure:"interval"`
	PairFilter PairFilter  `mapstructure:"pair_filter"`
	Retry      RetryConfig `mapstructure:"retry"`
}

// PairFilter represents trading pair filtering configuration
type PairFilter struct {
	QuoteAsset   string   `mapstructure:"quote_asset"`
	ExcludePairs []string `mapstructure:"exclude_pairs"`
}

// RetryConfig represents retry configuration
type RetryConfig struct {
	MaxAttempts       int           `mapstructure:"max_attempts"`
	Delay             time.Duration `mapstructure:"delay"`
	BackoffMultiplier float64       `mapstructure:"backoff_multiplier"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type  string      `mapstructure:"type"`
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
}

// MySQLConfig represents MySQL database configuration
type MySQLConfig struct {
	Host               string        `mapstructure:"host"`
	Port               int           `mapstructure:"port"`
	Database           string        `mapstructure:"database"`
	Username           string        `mapstructure:"username"`
	Password           string        `mapstructure:"password"`
	Charset            string        `mapstructure:"charset"`
	ParseTime          bool          `mapstructure:"parse_time"`
	MaxOpenConns       int           `mapstructure:"max_open_conns"`
	MaxIdleConns       int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime    time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime    time.Duration `mapstructure:"conn_max_idle_time"`
	SlowQueryThreshold time.Duration `mapstructure:"slow_query_threshold"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// StrategiesConfig represents all strategy configurations
type StrategiesConfig struct {
	Minority MinorityStrategy `mapstructure:"minority"`
	Whale    WhaleStrategy    `mapstructure:"whale"`
	Global   GlobalStrategy   `mapstructure:"global"`
}

// MinorityStrategy represents minority follower strategy configuration
type MinorityStrategy struct {
	Enabled                         bool    `mapstructure:"enabled"`
	Name                            string  `mapstructure:"name"`
	MinRatioDifference              float64 `mapstructure:"min_ratio_difference"`
	ConfirmationHours               int     `mapstructure:"confirmation_hours"`
	GenerateLongWhenShortRatioAbove float64 `mapstructure:"generate_long_when_short_ratio_above"`
	GenerateShortWhenLongRatioAbove float64 `mapstructure:"generate_short_when_long_ratio_above"`
	TrackingHours                   int     `mapstructure:"tracking_hours"`
	ProfitTargetPct                 float64 `mapstructure:"profit_target_pct"`
	StopLossPct                     float64 `mapstructure:"stop_loss_pct"`
}

// WhaleStrategy represents whale position analysis strategy configuration
type WhaleStrategy struct {
	Enabled                bool    `mapstructure:"enabled"`
	Name                   string  `mapstructure:"name"`
	MinRatioDifference     float64 `mapstructure:"min_ratio_difference"`
	WhalePositionThreshold float64 `mapstructure:"whale_position_threshold"`
	ConfirmationHours      int     `mapstructure:"confirmation_hours"`
	MinDivergence          float64 `mapstructure:"min_divergence"`
	TrackingHours          int     `mapstructure:"tracking_hours"`
	ProfitTargetPct        float64 `mapstructure:"profit_target_pct"`
	StopLossPct            float64 `mapstructure:"stop_loss_pct"`
}

// GlobalStrategy represents global strategy settings
type GlobalStrategy struct {
	MinVolume24h                float64 `mapstructure:"min_volume_24h"`
	MaxConcurrentSignalsPerPair int     `mapstructure:"max_concurrent_signals_per_pair"`
	SignalCooldownHours         int     `mapstructure:"signal_cooldown_hours"`
}

// StatisticsConfig represents statistics calculation configuration
type StatisticsConfig struct {
	CalculationInterval string                      `mapstructure:"calculation_interval"`
	Periods             []string                    `mapstructure:"periods"`
	Percentiles         []int                       `mapstructure:"percentiles"`
	Monitoring          StatisticsMonitoringConfig  `mapstructure:"monitoring"`
}

// StatisticsMonitoringConfig configures change detection thresholds
type StatisticsMonitoringConfig struct {
	Enabled                     bool    `mapstructure:"enabled"`
	WinRateChangeThreshold      float64 `mapstructure:"win_rate_change_threshold"`
	ProfitRatioChangeThreshold  float64 `mapstructure:"profit_ratio_change_threshold"`
	AvgProfitChangeThreshold    float64 `mapstructure:"avg_profit_change_threshold"`
	AvgLossChangeThreshold      float64 `mapstructure:"avg_loss_change_threshold"`
	ProfitFactorChangeThreshold float64 `mapstructure:"profit_factor_change_threshold"`
	SignalCountChangeThreshold  float64 `mapstructure:"signal_count_change_threshold"`
}

// NotificationsConfig represents all notification configurations
type NotificationsConfig struct {
	Telegram TelegramConfig `mapstructure:"telegram"`
	Email    EmailConfig    `mapstructure:"email"`
	Webhook  WebhookConfig  `mapstructure:"webhook"`
	Console  ConsoleConfig  `mapstructure:"console"`
}

// TelegramConfig represents Telegram notification configuration
type TelegramConfig struct {
	Enabled  bool     `mapstructure:"enabled"`
	BotToken string   `mapstructure:"bot_token"`
	ChatIDs  []string `mapstructure:"chat_ids"`
	Events   []string `mapstructure:"events"`
	Template string   `mapstructure:"template"`
}

// EmailConfig represents email notification configuration
type EmailConfig struct {
	Enabled  bool     `mapstructure:"enabled"`
	SMTPHost string   `mapstructure:"smtp_host"`
	SMTPPort int      `mapstructure:"smtp_port"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	From     string   `mapstructure:"from"`
	To       []string `mapstructure:"to"`
	Events   []string `mapstructure:"events"`
}

// WebhookConfig represents webhook notification configuration
type WebhookConfig struct {
	Enabled bool              `mapstructure:"enabled"`
	URL     string            `mapstructure:"url"`
	Method  string            `mapstructure:"method"`
	Headers map[string]string `mapstructure:"headers"`
	Timeout time.Duration     `mapstructure:"timeout"`
	Events  []string          `mapstructure:"events"`
}

// ConsoleConfig represents console notification configuration
type ConsoleConfig struct {
	Enabled bool     `mapstructure:"enabled"`
	Events  []string `mapstructure:"events"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string        `mapstructure:"level"`
	Format string        `mapstructure:"format"`
	Output []string      `mapstructure:"output"`
	File   FileLogConfig `mapstructure:"file"`
}

// FileLogConfig represents file logging configuration
type FileLogConfig struct {
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// MonitoringConfig represents monitoring configuration
type MonitoringConfig struct {
	Enabled     bool              `mapstructure:"enabled"`
	Metrics     MetricsConfig     `mapstructure:"metrics"`
	HealthCheck HealthCheckConfig `mapstructure:"health_check"`
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Path    string `mapstructure:"path"`
}

// HealthCheckConfig represents health check configuration
type HealthCheckConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

// FeaturesConfig represents feature flags
type FeaturesConfig struct {
	BacktestMode bool `mapstructure:"backtest_mode"`
	DryRun       bool `mapstructure:"dry_run"`
	DebugSignals bool `mapstructure:"debug_signals"`
}
