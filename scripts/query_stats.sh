#!/bin/bash

# 胜率查询工具
# 使用方法:
#   ./scripts/query_stats.sh                 # 查询所有策略的整体胜率
#   ./scripts/query_stats.sh minority 24h    # 查询Minority策略24小时胜率
#   ./scripts/query_stats.sh whale 7d        # 查询Whale策略7天胜率

# 数据库配置
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="futures_analysis"
DB_USER="root"
DB_PASS="123456"

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

function query_all_stats() {
    echo -e "${GREEN}=== 所有策略的最新整体胜率 ===${NC}\n"

    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -t <<EOF
SELECT
    strategy_name AS '策略',
    period_label AS '周期',
    total_signals AS '信号数',
    confirmed_signals AS '确认',
    profitable_signals AS '盈利',
    losing_signals AS '亏损',
    neutral_signals AS '中性',
    COALESCE(win_rate, 0) AS '胜率(%)',
    COALESCE(kline_close_win_rate, 0) AS 'K线胜率(%)',
    total_kline_hours AS 'K线时数',
    DATE_FORMAT(calculated_at, '%m-%d %H:%i') AS '计算时间'
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
EOF
}

function query_strategy_stats() {
    local strategy=$1
    local period=${2:-"24h"}

    # 转换策略名称
    case $strategy in
        minority|min)
            strategy_full="Minority Follower"
            ;;
        whale|w)
            strategy_full="Whale Position Analysis"
            ;;
        *)
            strategy_full=$strategy
            ;;
    esac

    echo -e "${BLUE}=== ${strategy_full} 策略 - ${period} 周期统计 ===${NC}\n"

    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -t <<EOF
SELECT
    COALESCE(symbol, '【整体】') AS '交易对',
    total_signals AS '信号',
    confirmed_signals AS '确认',
    profitable_signals AS '盈利',
    losing_signals AS '亏损',
    neutral_signals AS '中性',
    ROUND(COALESCE(win_rate, 0), 2) AS '胜率(%)',
    ROUND(COALESCE(kline_theoretical_win_rate, 0), 2) AS '理论胜率(%)',
    ROUND(COALESCE(kline_close_win_rate, 0), 2) AS 'K线胜率(%)',
    total_kline_hours AS 'K线时数',
    ROUND(COALESCE(avg_hourly_return_pct, 0), 2) AS '平均小时收益(%)'
FROM strategy_statistics
WHERE strategy_name = '$strategy_full'
    AND period_label = '$period'
    AND calculated_at = (
        SELECT MAX(calculated_at)
        FROM strategy_statistics s2
        WHERE s2.strategy_name = strategy_statistics.strategy_name
            AND s2.period_label = strategy_statistics.period_label
            AND COALESCE(s2.symbol, '') = COALESCE(strategy_statistics.symbol, '')
    )
ORDER BY
    CASE WHEN symbol IS NULL THEN 0 ELSE 1 END,
    total_signals DESC;
EOF
}

function query_top_performers() {
    echo -e "${YELLOW}=== 胜率最高的交易对 (Top 10) ===${NC}\n"

    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -t <<EOF
SELECT
    strategy_name AS '策略',
    symbol AS '交易对',
    period_label AS '周期',
    ROUND(win_rate, 2) AS '胜率(%)',
    total_signals AS '信号数',
    profitable_signals AS '盈利',
    losing_signals AS '亏损',
    ROUND(COALESCE(kline_close_win_rate, 0), 2) AS 'K线胜率(%)'
FROM strategy_statistics
WHERE symbol IS NOT NULL
    AND win_rate IS NOT NULL
    AND (profitable_signals + losing_signals + neutral_signals) >= 3
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
EOF
}

function show_help() {
    echo -e "${GREEN}胜率查询工具使用说明${NC}"
    echo ""
    echo "用法:"
    echo "  $0                      # 查询所有策略的整体胜率"
    echo "  $0 minority [period]    # 查询Minority策略（可指定周期: 24h/7d/30d/all）"
    echo "  $0 whale [period]       # 查询Whale策略"
    echo "  $0 top                  # 查询胜率最高的交易对"
    echo ""
    echo "示例:"
    echo "  $0                      # 查询所有"
    echo "  $0 minority             # 查询Minority策略24h数据"
    echo "  $0 minority 7d          # 查询Minority策略7天数据"
    echo "  $0 whale 30d            # 查询Whale策略30天数据"
    echo "  $0 top                  # 查询表现最佳的交易对"
    echo ""
}

# 主程序
case ${1:-all} in
    all|"")
        query_all_stats
        ;;
    minority|min)
        query_strategy_stats "$1" "$2"
        ;;
    whale|w)
        query_strategy_stats "$1" "$2"
        ;;
    top)
        query_top_performers
        ;;
    help|-h|--help)
        show_help
        ;;
    *)
        echo "未知命令: $1"
        show_help
        exit 1
        ;;
esac
