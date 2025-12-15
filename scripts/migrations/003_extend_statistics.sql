-- Migration: 003_extend_statistics.sql
-- Description: Extend strategy_statistics table with kline-based win rate metrics
-- Author: Claude Code
-- Date: 2025-12-13

-- Add kline-based win rate and performance metrics to strategy_statistics table

ALTER TABLE strategy_statistics
    -- Kline win rate metrics
    ADD COLUMN kline_theoretical_win_rate DECIMAL(10,4) DEFAULT NULL COMMENT '理论胜率 - based on high price (%)',
    ADD COLUMN kline_close_win_rate DECIMAL(10,4) DEFAULT NULL COMMENT '小时K线胜率 - based on close price (%)',
    ADD COLUMN total_kline_hours INT DEFAULT 0 COMMENT '总K线小时数',
    ADD COLUMN profitable_kline_hours_high INT DEFAULT 0 COMMENT '最高价盈利的小时数',
    ADD COLUMN profitable_kline_hours_close INT DEFAULT 0 COMMENT '收盘价盈利的小时数',

    -- Hourly return statistics
    ADD COLUMN avg_hourly_return_pct DECIMAL(10,4) DEFAULT NULL COMMENT '平均小时收益率',
    ADD COLUMN max_hourly_return_pct DECIMAL(10,4) DEFAULT NULL COMMENT '最大小时收益率',
    ADD COLUMN min_hourly_return_pct DECIMAL(10,4) DEFAULT NULL COMMENT '最小小时收益率',

    -- Theoretical maximum profit/loss
    ADD COLUMN avg_max_potential_profit_pct DECIMAL(10,4) DEFAULT NULL COMMENT '平均最高价理论收益',
    ADD COLUMN avg_max_potential_loss_pct DECIMAL(10,4) DEFAULT NULL COMMENT '平均最低价最大回撤',

    -- Index for efficient queries on kline win rates
    ADD INDEX idx_kline_win_rates (kline_theoretical_win_rate, kline_close_win_rate);
