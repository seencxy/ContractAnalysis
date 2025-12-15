-- Migration: 002_add_kline_tracking.sql
-- Description: Create signal_kline_tracking table for hourly kline-based tracking
-- Author: Claude Code
-- Date: 2025-12-13

-- Table: signal_kline_tracking
-- Purpose: Store hourly kline tracking data for signals
-- Relationship: One signal can have many kline tracking records (one per hour)

CREATE TABLE IF NOT EXISTS signal_kline_tracking (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    signal_id VARCHAR(36) NOT NULL COMMENT 'Associated signal ID (UUID)',

    -- K-line time information
    kline_open_time TIMESTAMP NOT NULL COMMENT 'Kline open time (整点)',
    kline_close_time TIMESTAMP NOT NULL COMMENT 'Kline close time',
    hours_since_signal DECIMAL(10,2) NOT NULL COMMENT 'Hours elapsed since signal generated',

    -- OHLCV data
    open_price DECIMAL(20,8) NOT NULL COMMENT 'Open price',
    high_price DECIMAL(20,8) NOT NULL COMMENT 'High price',
    low_price DECIMAL(20,8) NOT NULL COMMENT 'Low price',
    close_price DECIMAL(20,8) NOT NULL COMMENT 'Close price',
    volume DECIMAL(20,4) NOT NULL COMMENT 'Volume',
    quote_volume DECIMAL(20,4) NOT NULL COMMENT 'Quote asset volume',

    -- Price change percentages (relative to signal price)
    open_change_pct DECIMAL(10,4) NOT NULL COMMENT 'Open price change %',
    high_change_pct DECIMAL(10,4) NOT NULL COMMENT 'High price change %',
    low_change_pct DECIMAL(10,4) NOT NULL COMMENT 'Low price change %',
    close_change_pct DECIMAL(10,4) NOT NULL COMMENT 'Close price change %',

    -- Hourly return
    hourly_return_pct DECIMAL(10,4) NOT NULL COMMENT 'Hourly return: (close-open)/open*100',

    -- Theoretical maximum profit/loss
    max_potential_profit_pct DECIMAL(10,4) NOT NULL COMMENT 'Max profit at high price',
    max_potential_loss_pct DECIMAL(10,4) NOT NULL COMMENT 'Max drawdown at low price',

    -- Profitability flags
    is_profitable_at_high BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Is profitable at high price',
    is_profitable_at_close BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Is profitable at close price',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Indexes for efficient queries
    INDEX idx_signal_id (signal_id),
    INDEX idx_kline_open_time (kline_open_time),
    INDEX idx_signal_kline (signal_id, kline_open_time),
    INDEX idx_profitable_high (is_profitable_at_high),
    INDEX idx_profitable_close (is_profitable_at_close),

    -- Foreign key constraint
    FOREIGN KEY (signal_id) REFERENCES signals(signal_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='Signal kline tracking table - tracks hourly kline data for each signal';
