package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Load loads the configuration from the specified file
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Default config file name and paths
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("$HOME/.contractanalysis")
	}

	// Enable environment variable override
	v.SetEnvPrefix("CA") // ContractAnalysis
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Set defaults
	setDefaults(v)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, use defaults
			fmt.Fprintf(os.Stderr, "Warning: Config file not found, using defaults\n")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config into struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "Binance Futures Analysis")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.timezone", "UTC")

	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")

	// Binance defaults
	v.SetDefault("binance.api_url", "https://fapi.binance.com")
	v.SetDefault("binance.rate_limit.requests_per_minute", 1200)
	v.SetDefault("binance.rate_limit.weight_per_minute", 2400)
	v.SetDefault("binance.timeout", "10s")

	// Collection defaults
	v.SetDefault("collection.enabled", true)
	v.SetDefault("collection.interval", "0 * * * *")
	v.SetDefault("collection.pair_filter.quote_asset", "USDT")
	v.SetDefault("collection.retry.max_attempts", 3)
	v.SetDefault("collection.retry.delay", "5s")
	v.SetDefault("collection.retry.backoff_multiplier", 2.0)

	// Database defaults
	v.SetDefault("database.type", "mysql")
	v.SetDefault("database.mysql.host", "localhost")
	v.SetDefault("database.mysql.port", 3306)
	v.SetDefault("database.mysql.database", "futures_analysis")
	v.SetDefault("database.mysql.charset", "utf8mb4")
	v.SetDefault("database.mysql.parse_time", true)
	v.SetDefault("database.mysql.max_open_conns", 25)
	v.SetDefault("database.mysql.max_idle_conns", 5)
	v.SetDefault("database.mysql.conn_max_lifetime", "5m")
	v.SetDefault("database.mysql.conn_max_idle_time", "10m")
	v.SetDefault("database.mysql.slow_query_threshold", "200ms")

	v.SetDefault("database.redis.host", "localhost")
	v.SetDefault("database.redis.port", 6379)
	v.SetDefault("database.redis.db", 0)
	v.SetDefault("database.redis.pool_size", 10)
	v.SetDefault("database.redis.min_idle_conns", 5)
	v.SetDefault("database.redis.max_retries", 3)
	v.SetDefault("database.redis.dial_timeout", "5s")
	v.SetDefault("database.redis.read_timeout", "3s")
	v.SetDefault("database.redis.write_timeout", "3s")

	// Strategy defaults
	v.SetDefault("strategies.minority.enabled", true)
	v.SetDefault("strategies.minority.name", "Minority Follower")
	v.SetDefault("strategies.minority.min_ratio_difference", 75.0)
	v.SetDefault("strategies.minority.confirmation_hours", 2)
	v.SetDefault("strategies.minority.generate_long_when_short_ratio_above", 75.0)
	v.SetDefault("strategies.minority.generate_short_when_long_ratio_above", 75.0)
	v.SetDefault("strategies.minority.tracking_hours", 24)
	v.SetDefault("strategies.minority.profit_target_pct", 5.0)
	v.SetDefault("strategies.minority.stop_loss_pct", 2.0)

	v.SetDefault("strategies.whale.enabled", true)
	v.SetDefault("strategies.whale.name", "Whale Position Analysis")
	v.SetDefault("strategies.whale.min_ratio_difference", 75.0)
	v.SetDefault("strategies.whale.whale_position_threshold", 60.0)
	v.SetDefault("strategies.whale.confirmation_hours", 2)
	v.SetDefault("strategies.whale.min_divergence", 20.0)
	v.SetDefault("strategies.whale.tracking_hours", 24)
	v.SetDefault("strategies.whale.profit_target_pct", 5.0)
	v.SetDefault("strategies.whale.stop_loss_pct", 2.0)

	v.SetDefault("strategies.global.min_volume_24h", 1000000)
	v.SetDefault("strategies.global.max_concurrent_signals_per_pair", 3)
	v.SetDefault("strategies.global.signal_cooldown_hours", 6)

	// Statistics defaults
	v.SetDefault("statistics.calculation_interval", "0 */6 * * *")
	v.SetDefault("statistics.periods", []string{"24h", "7d", "30d", "all"})
	v.SetDefault("statistics.percentiles", []int{25, 50, 75, 90, 95})

	// Notification defaults
	v.SetDefault("notifications.console.enabled", true)
	v.SetDefault("notifications.console.events", []string{"signal_generated", "signal_confirmed", "signal_invalidated", "signal_outcome"})

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output", []string{"stdout"})
	v.SetDefault("logging.file.path", "logs/app.log")
	v.SetDefault("logging.file.max_size", 100)
	v.SetDefault("logging.file.max_backups", 10)
	v.SetDefault("logging.file.max_age", 30)
	v.SetDefault("logging.file.compress", true)

	// Monitoring defaults
	v.SetDefault("monitoring.enabled", true)
	v.SetDefault("monitoring.metrics.enabled", true)
	v.SetDefault("monitoring.metrics.port", 9090)
	v.SetDefault("monitoring.metrics.path", "/metrics")
	v.SetDefault("monitoring.health_check.enabled", true)
	v.SetDefault("monitoring.health_check.path", "/health")

	// Feature flags defaults
	v.SetDefault("features.backtest_mode", false)
	v.SetDefault("features.dry_run", false)
	v.SetDefault("features.debug_signals", false)
}

// validate validates the configuration
func validate(config *Config) error {
	// Validate app
	if config.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}

	// Validate Binance config if collection is enabled
	if config.Collection.Enabled {
		if config.Binance.APIURL == "" {
			return fmt.Errorf("binance.api_url is required when collection is enabled")
		}
	}

	// Validate database
	if config.Database.Type != "mysql" && config.Database.Type != "redis" {
		return fmt.Errorf("database.type must be 'mysql' or 'redis', got: %s", config.Database.Type)
	}

	if config.Database.Type == "mysql" {
		if config.Database.MySQL.Host == "" {
			return fmt.Errorf("database.mysql.host is required")
		}
		if config.Database.MySQL.Database == "" {
			return fmt.Errorf("database.mysql.database is required")
		}
	}

	// Validate strategies
	if config.Strategies.Minority.Enabled {
		if config.Strategies.Minority.MinRatioDifference < 50 || config.Strategies.Minority.MinRatioDifference > 100 {
			return fmt.Errorf("strategies.minority.min_ratio_difference must be between 50 and 100")
		}
	}

	if config.Strategies.Whale.Enabled {
		if config.Strategies.Whale.WhalePositionThreshold < 0 || config.Strategies.Whale.WhalePositionThreshold > 100 {
			return fmt.Errorf("strategies.whale.whale_position_threshold must be between 0 and 100")
		}
	}

	// Validate logging
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[config.Logging.Level] {
		return fmt.Errorf("logging.level must be one of: debug, info, warn, error")
	}

	validLogFormats := map[string]bool{"json": true, "console": true}
	if !validLogFormats[config.Logging.Format] {
		return fmt.Errorf("logging.format must be one of: json, console")
	}

	return nil
}
