-- 查询胜率统计
-- 使用方法: mysql -u root -p futures_analysis < scripts/query_win_rate.sql

-- 1. 查询所有策略的最新整体胜率（不分交易对）
SELECT
    strategy_name AS '策略名称',
    period_label AS '时间周期',
    total_signals AS '总信号数',
    confirmed_signals AS '已确认',
    profitable_signals AS '盈利',
    losing_signals AS '亏损',
    neutral_signals AS '中性',
    COALESCE(win_rate, 0) AS '胜率(%)',
    COALESCE(kline_theoretical_win_rate, 0) AS 'K线理论胜率(%)',
    COALESCE(kline_close_win_rate, 0) AS 'K线实际胜率(%)',
    total_kline_hours AS 'K线小时数',
    DATE_FORMAT(calculated_at, '%Y-%m-%d %H:%i') AS '计算时间'
FROM strategy_statistics
WHERE symbol IS NULL
    AND calculated_at = (
        SELECT MAX(calculated_at)
        FROM strategy_statistics s2
        WHERE s2.strategy_name = strategy_statistics.strategy_name
            AND s2.period_label = strategy_statistics.period_label
            AND s2.symbol IS NULL
    )
ORDER BY strategy_name,
    FIELD(period_label, '24h', '7d', '30d', 'all');

-- 2. 查询特定策略的按交易对分组的胜率
SELECT
    strategy_name AS '策略名称',
    COALESCE(symbol, '整体') AS '交易对',
    period_label AS '周期',
    total_signals AS '信号数',
    COALESCE(win_rate, 0) AS '胜率(%)',
    COALESCE(kline_close_win_rate, 0) AS 'K线胜率(%)',
    profitable_signals AS '盈利',
    losing_signals AS '亏损',
    neutral_signals AS '中性'
FROM strategy_statistics
WHERE strategy_name = 'Minority Follower'
    AND period_label = '24h'
    AND calculated_at = (
        SELECT MAX(calculated_at)
        FROM strategy_statistics s2
        WHERE s2.strategy_name = strategy_statistics.strategy_name
            AND s2.period_label = strategy_statistics.period_label
            AND COALESCE(s2.symbol, '') = COALESCE(strategy_statistics.symbol, '')
    )
ORDER BY total_signals DESC;

-- 3. 查询胜率最高的交易对（需要有实际数据）
SELECT
    strategy_name AS '策略名称',
    symbol AS '交易对',
    period_label AS '周期',
    COALESCE(win_rate, 0) AS '胜率(%)',
    total_signals AS '信号数',
    profitable_signals AS '盈利',
    losing_signals AS '亏损'
FROM strategy_statistics
WHERE symbol IS NOT NULL
    AND win_rate IS NOT NULL
    AND period_label = '24h'
    AND calculated_at = (
        SELECT MAX(calculated_at)
        FROM strategy_statistics s2
        WHERE s2.strategy_name = strategy_statistics.strategy_name
            AND s2.period_label = strategy_statistics.period_label
            AND s2.symbol = strategy_statistics.symbol
    )
ORDER BY win_rate DESC
LIMIT 10;
