#!/bin/bash

echo "=== End-to-End Integration Verification ==="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Track pass/fail
PASSED=0
FAILED=0

# Check 1: Database migrations
echo "--- Check 1: Database Tables ---"
TABLES=$(docker exec mysql mysql -u root -p123456 -D futures_analysis -e "SHOW TABLES LIKE 'signal_kline_tracking';" 2>&1 | grep -v "Using a password" | grep signal_kline_tracking)
if [ -n "$TABLES" ]; then
    echo -e "${GREEN}✓${NC} signal_kline_tracking table exists"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} signal_kline_tracking table not found"
    ((FAILED++))
fi

COLUMNS=$(docker exec mysql mysql -u root -p123456 -D futures_analysis -e "DESCRIBE strategy_statistics;" 2>&1 | grep -v "Using a password" | grep kline_theoretical_win_rate)
if [ -n "$COLUMNS" ]; then
    echo -e "${GREEN}✓${NC} strategy_statistics extended with kline fields"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} strategy_statistics missing kline fields"
    ((FAILED++))
fi

# Check 2: Build success
echo ""
echo "--- Check 2: Application Build ---"
if [ -f "/tmp/futures-analysis" ]; then
    echo -e "${GREEN}✓${NC} Application binary built successfully"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Application binary not found"
    ((FAILED++))
fi

# Check 3: Verify main components in code
echo ""
echo "--- Check 3: Code Integration Points ---"

# Check binance client has kline methods
if grep -q "GetKlines" /Volumes/tiam/Project/ContractAnalysis/internal/infrastructure/binance/client.go; then
    echo -e "${GREEN}✓${NC} Binance client has GetKlines method"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Binance client missing GetKlines method"
    ((FAILED++))
fi

# Check tracker has TrackAllKlines
if grep -q "TrackAllKlines" /Volumes/tiam/Project/ContractAnalysis/internal/usecase/tracker.go; then
    echo -e "${GREEN}✓${NC} Tracker has TrackAllKlines method"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Tracker missing TrackAllKlines method"
    ((FAILED++))
fi

# Check statistics calculator has kline metrics
if grep -q "calculateKlineMetrics" /Volumes/tiam/Project/ContractAnalysis/internal/usecase/statistics_calculator.go; then
    echo -e "${GREEN}✓${NC} Statistics calculator has calculateKlineMetrics method"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Statistics calculator missing calculateKlineMetrics method"
    ((FAILED++))
fi

# Check scheduler has kline tracking job
if grep -q "AddKlineTrackingJob" /Volumes/tiam/Project/ContractAnalysis/internal/infrastructure/scheduler/cron.go; then
    echo -e "${GREEN}✓${NC} Scheduler has AddKlineTrackingJob method"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Scheduler missing AddKlineTrackingJob method"
    ((FAILED++))
fi

# Check main.go registers the job
if grep -q "AddKlineTrackingJob" /Volumes/tiam/Project/ContractAnalysis/main.go; then
    echo -e "${GREEN}✓${NC} main.go registers kline tracking job"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} main.go doesn't register kline tracking job"
    ((FAILED++))
fi

# Check 4: Repository methods
echo ""
echo "--- Check 4: Repository Implementation ---"

if grep -q "CreateKlineTracking" /Volumes/tiam/Project/ContractAnalysis/internal/infrastructure/persistence/mysql/signal.go; then
    echo -e "${GREEN}✓${NC} Repository has CreateKlineTracking method"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Repository missing CreateKlineTracking method"
    ((FAILED++))
fi

if grep -q "GetKlineTrackingBySignal" /Volumes/tiam/Project/ContractAnalysis/internal/infrastructure/persistence/mysql/signal.go; then
    echo -e "${GREEN}✓${NC} Repository has GetKlineTrackingBySignal method"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Repository missing GetKlineTrackingBySignal method"
    ((FAILED++))
fi

# Check 5: Entity definitions
echo ""
echo "--- Check 5: Entity Definitions ---"

if [ -f "/Volumes/tiam/Project/ContractAnalysis/internal/domain/entity/signal_kline_tracking.go" ]; then
    echo -e "${GREEN}✓${NC} SignalKlineTracking entity exists"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} SignalKlineTracking entity not found"
    ((FAILED++))
fi

if grep -q "Kline struct" /Volumes/tiam/Project/ContractAnalysis/internal/domain/entity/signal_kline_tracking.go; then
    echo -e "${GREEN}✓${NC} Kline struct defined"
    ((PASSED++))
else
    echo -e "${RED}✗${NC} Kline struct not found"
    ((FAILED++))
fi

# Summary
echo ""
echo "========================================="
echo "Test Results:"
echo -e "  ${GREEN}Passed: $PASSED${NC}"
echo -e "  ${RED}Failed: $FAILED${NC}"
echo "========================================="

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All integration checks passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some integration checks failed!${NC}"
    exit 1
fi
