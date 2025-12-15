-- 001_initial_schema.sql
-- Initial database schema for Binance Futures Analysis System

-- Table: trading_pairs
-- Stores all USDT-margined futures pairs
CREATE TABLE IF NOT EXISTS trading_pairs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(50) NOT NULL UNIQUE COMMENT 'Trading pair symbol, e.g., BTCUSDT',
    is_active BOOLEAN DEFAULT TRUE COMMENT 'Whether the pair is actively traded',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_symbol (symbol),
    INDEX idx_is_active (is_active),
    INDEX idx_symbol_active (symbol, is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Trading pairs table';

-- Table: market_data
-- Stores collected long/short ratio data
CREATE TABLE IF NOT EXISTS market_data (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(50) NOT NULL COMMENT 'Trading pair symbol',
    timestamp TIMESTAMP NOT NULL COMMENT 'Data collection timestamp',

    -- Long/Short ratio by account count
    long_account_ratio DECIMAL(10,4) NOT NULL COMMENT 'Long account ratio (percentage)',
    short_account_ratio DECIMAL(10,4) NOT NULL COMMENT 'Short account ratio (percentage)',

    -- Long/Short ratio by position size (whale positions)
    long_position_ratio DECIMAL(10,4) NOT NULL COMMENT 'Long position ratio by size (percentage)',
    short_position_ratio DECIMAL(10,4) NOT NULL COMMENT 'Short position ratio by size (percentage)',

    -- Trader count ratio
    long_trader_count INT NOT NULL COMMENT 'Number of long traders',
    short_trader_count INT NOT NULL COMMENT 'Number of short traders',

    -- Price at collection time
    price DECIMAL(20,8) NOT NULL COMMENT 'Price at data collection time',

    -- 24h volume (optional, for filtering)
    volume_24h DECIMAL(20,2) DEFAULT NULL COMMENT '24-hour trading volume',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_symbol_timestamp (symbol, timestamp),
    INDEX idx_symbol (symbol),
    INDEX idx_timestamp (timestamp),
    INDEX idx_symbol_timestamp (symbol, timestamp),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Market data table storing long/short ratios';

-- Table: signals
-- Stores generated trading signals
CREATE TABLE IF NOT EXISTS signals (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    signal_id VARCHAR(36) NOT NULL UNIQUE COMMENT 'UUID of the signal',
    symbol VARCHAR(50) NOT NULL COMMENT 'Trading pair symbol',

    -- Signal details
    signal_type VARCHAR(20) NOT NULL COMMENT 'LONG or SHORT',
    strategy_name VARCHAR(50) NOT NULL COMMENT 'Strategy that generated the signal',

    -- Market conditions at signal time
    generated_at TIMESTAMP NOT NULL COMMENT 'Signal generation timestamp',
    price_at_signal DECIMAL(20,8) NOT NULL COMMENT 'Price when signal was generated',
    long_account_ratio DECIMAL(10,4) NOT NULL,
    short_account_ratio DECIMAL(10,4) NOT NULL,
    long_position_ratio DECIMAL(10,4) NOT NULL,
    short_position_ratio DECIMAL(10,4) NOT NULL,
    long_trader_count INT NOT NULL,
    short_trader_count INT NOT NULL,

    -- Confirmation tracking
    confirmation_start TIMESTAMP NOT NULL COMMENT 'Start of confirmation period',
    confirmation_end TIMESTAMP NOT NULL COMMENT 'End of confirmation period',
    is_confirmed BOOLEAN DEFAULT FALSE,
    confirmed_at TIMESTAMP NULL,

    -- Signal status
    status VARCHAR(20) NOT NULL COMMENT 'PENDING, CONFIRMED, INVALIDATED, TRACKING, CLOSED',

    -- Metadata
    reason TEXT COMMENT 'Reason why signal was generated',
    config_snapshot JSON COMMENT 'Strategy configuration at signal time',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_signal_id (signal_id),
    INDEX idx_symbol (symbol),
    INDEX idx_status (status),
    INDEX idx_generated_at (generated_at),
    INDEX idx_strategy (strategy_name),
    INDEX idx_symbol_status (symbol, status),
    INDEX idx_status_generated (status, generated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Trading signals table';

-- Table: signal_tracking
-- Tracks price movements after signals
CREATE TABLE IF NOT EXISTS signal_tracking (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    signal_id VARCHAR(36) NOT NULL COMMENT 'Reference to signal',

    -- Tracking details
    tracked_at TIMESTAMP NOT NULL COMMENT 'Tracking timestamp',
    hours_elapsed DECIMAL(10,2) NOT NULL COMMENT 'Hours since signal generation',

    current_price DECIMAL(20,8) NOT NULL,
    price_change_pct DECIMAL(10,4) NOT NULL COMMENT 'Price change percentage',

    -- Peak/trough tracking
    highest_price DECIMAL(20,8) NOT NULL,
    highest_price_pct DECIMAL(10,4) NOT NULL,
    highest_price_at TIMESTAMP NOT NULL,

    lowest_price DECIMAL(20,8) NOT NULL,
    lowest_price_pct DECIMAL(10,4) NOT NULL,
    lowest_price_at TIMESTAMP NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_signal_id (signal_id),
    INDEX idx_tracked_at (tracked_at),
    INDEX idx_signal_tracked (signal_id, tracked_at),

    FOREIGN KEY (signal_id) REFERENCES signals(signal_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Signal tracking table';

-- Table: signal_outcomes
-- Final outcomes of signals for statistics
CREATE TABLE IF NOT EXISTS signal_outcomes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    signal_id VARCHAR(36) NOT NULL UNIQUE COMMENT 'Reference to signal',

    -- Outcome details
    outcome VARCHAR(20) NOT NULL COMMENT 'PROFIT, LOSS, NEUTRAL, TIMEOUT',

    -- Performance metrics
    max_favorable_move_pct DECIMAL(10,4) NOT NULL COMMENT 'Maximum favorable price move',
    max_adverse_move_pct DECIMAL(10,4) NOT NULL COMMENT 'Maximum adverse price move',
    final_price_change_pct DECIMAL(10,4) NOT NULL COMMENT 'Final price change',

    -- Timing
    hours_to_peak INT COMMENT 'Hours to reach peak favorable move',
    hours_to_trough INT COMMENT 'Hours to reach trough adverse move',
    total_tracking_hours INT NOT NULL,

    -- Additional metrics
    profit_target_hit BOOLEAN DEFAULT FALSE,
    stop_loss_hit BOOLEAN DEFAULT FALSE,

    closed_at TIMESTAMP NOT NULL COMMENT 'When signal was closed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_signal_id (signal_id),
    INDEX idx_outcome (outcome),
    INDEX idx_closed_at (closed_at),

    FOREIGN KEY (signal_id) REFERENCES signals(signal_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Signal outcomes table for statistics';

-- Table: strategy_statistics
-- Aggregated strategy performance statistics
CREATE TABLE IF NOT EXISTS strategy_statistics (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    strategy_name VARCHAR(50) NOT NULL,
    symbol VARCHAR(50) NULL COMMENT 'NULL for overall stats',

    -- Time period
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    period_label VARCHAR(20) NOT NULL COMMENT '24h, 7d, 30d, all',

    -- Signal counts
    total_signals INT NOT NULL DEFAULT 0,
    confirmed_signals INT NOT NULL DEFAULT 0,
    invalidated_signals INT NOT NULL DEFAULT 0,

    -- Outcome counts
    profitable_signals INT NOT NULL DEFAULT 0,
    losing_signals INT NOT NULL DEFAULT 0,
    neutral_signals INT NOT NULL DEFAULT 0,

    -- Performance metrics
    win_rate DECIMAL(10,4) COMMENT 'Win rate percentage',
    avg_profit_pct DECIMAL(10,4) COMMENT 'Average profit percentage',
    avg_loss_pct DECIMAL(10,4) COMMENT 'Average loss percentage',
    avg_holding_hours DECIMAL(10,2) COMMENT 'Average holding time in hours',

    -- Best/Worst
    best_signal_pct DECIMAL(10,4),
    worst_signal_pct DECIMAL(10,4),

    -- Profit factor
    profit_factor DECIMAL(10,4) COMMENT 'Gross profit / Gross loss',

    calculated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_strategy_symbol_period (strategy_name, symbol, period_label, period_start),
    INDEX idx_strategy (strategy_name),
    INDEX idx_symbol (symbol),
    INDEX idx_period (period_label),
    INDEX idx_period_range (period_start, period_end),
    INDEX idx_calculated_at (calculated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Strategy statistics table';

-- Table: notifications
-- Log of sent notifications
CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    signal_id VARCHAR(36) NULL COMMENT 'Related signal ID if applicable',

    notification_type VARCHAR(50) NOT NULL COMMENT 'Type of notification',
    channel VARCHAR(50) NOT NULL COMMENT 'telegram, email, webhook, console',

    recipient VARCHAR(255) NOT NULL COMMENT 'Recipient identifier',
    message TEXT NOT NULL COMMENT 'Notification message',

    status VARCHAR(20) NOT NULL COMMENT 'SENT, FAILED, PENDING',
    error_message TEXT NULL COMMENT 'Error message if failed',

    sent_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_signal_id (signal_id),
    INDEX idx_channel (channel),
    INDEX idx_status (status),
    INDEX idx_sent_at (sent_at),
    INDEX idx_type_channel (notification_type, channel),

    FOREIGN KEY (signal_id) REFERENCES signals(signal_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Notifications log table';
