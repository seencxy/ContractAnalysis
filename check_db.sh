#!/bin/bash
echo "=== Checking Database Status ==="
echo ""

# Check if we can connect
if ! command -v mysql &> /dev/null; then
    echo "MySQL client not found. Please run these queries manually:"
    echo ""
    echo "-- Check table counts:"
    echo "SELECT 'market_data' as tbl, COUNT(*) as cnt FROM market_data;"
    echo "SELECT 'signals' as tbl, COUNT(*) as cnt FROM signals;"  
    echo ""
    echo "-- Check recent market data:"
    echo "SELECT symbol, long_account_ratio, short_account_ratio, timestamp"
    echo "FROM market_data ORDER BY timestamp DESC LIMIT 10;"
    echo ""
    echo "-- Check if any signals exist:"
    echo "SELECT signal_id, symbol, signal_type, status, generated_at FROM signals LIMIT 10;"
    exit 0
fi

# Run checks
mysql -u root -p123456 futures_analysis << 'SQL'
SELECT '=== Table Counts ===' as info;
SELECT 'market_data' as table_name, COUNT(*) as count FROM market_data
UNION ALL SELECT 'signals', COUNT(*) FROM signals;

SELECT '\n=== Recent Market Data (Top 5) ===' as info;
SELECT symbol, long_account_ratio, short_account_ratio, 
       ROUND(long_account_ratio - short_account_ratio, 2) as diff,
       timestamp
FROM market_data 
ORDER BY timestamp DESC LIMIT 5;

SELECT '\n=== Check for Extreme Ratios (>= 75%) ===' as info;
SELECT symbol, long_account_ratio, short_account_ratio, timestamp
FROM market_data
WHERE long_account_ratio >= 75 OR short_account_ratio >= 75
ORDER BY timestamp DESC LIMIT 10;

SELECT '\n=== Signals Status ===' as info;
SELECT COUNT(*) as total,
       SUM(CASE WHEN status='PENDING' THEN 1 ELSE 0 END) as pending,
       SUM(CASE WHEN status='CONFIRMED' THEN 1 ELSE 0 END) as confirmed
FROM signals;
SQL
