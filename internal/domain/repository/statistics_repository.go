package repository

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// StrategyStatistics represents aggregated statistics for a strategy
type StrategyStatistics struct {
	ID           int64
	StrategyName string
	Symbol       *string // nil for overall stats
	PeriodStart  time.Time
	PeriodEnd    time.Time
	PeriodLabel  string // "24h", "7d", "30d", "all"

	// Signal counts
	TotalSignals       int
	ConfirmedSignals   int
	InvalidatedSignals int

	// Outcome counts
	ProfitableSignals int
	LosingSignals     int
	NeutralSignals    int

	// Performance metrics
	WinRate         *decimal.Decimal
	AvgProfitPct    *decimal.Decimal
	AvgLossPct      *decimal.Decimal
	AvgHoldingHours *decimal.Decimal

	// Best/Worst
	BestSignalPct  *decimal.Decimal
	WorstSignalPct *decimal.Decimal

	// Profit factor
	ProfitFactor *decimal.Decimal

	// Kline-based win rate metrics
	KlineTheoreticalWinRate   *decimal.Decimal // Win rate based on high price
	KlineCloseWinRate         *decimal.Decimal // Win rate based on close price
	TotalKlineHours           int              // Total kline hours tracked
	ProfitableKlineHoursHigh  int              // Hours profitable at high price
	ProfitableKlineHoursClose int              // Hours profitable at close price

	// Hourly return statistics
	AvgHourlyReturnPct *decimal.Decimal // Average hourly return
	MaxHourlyReturnPct *decimal.Decimal // Maximum hourly return
	MinHourlyReturnPct *decimal.Decimal // Minimum hourly return

	// Theoretical maximum profit/loss
	AvgMaxPotentialProfitPct *decimal.Decimal // Average max potential profit at high
	AvgMaxPotentialLossPct   *decimal.Decimal // Average max drawdown at low

	CalculatedAt time.Time
}

// StatisticsRepository defines the interface for statistics storage
type StatisticsRepository interface {
	// Create creates a new statistics record
	Create(ctx context.Context, stats *StrategyStatistics) error

	// CreateOrUpdate creates or updates a statistics record
	CreateOrUpdate(ctx context.Context, stats *StrategyStatistics) error

	// GetByStrategyAndPeriod retrieves statistics for a strategy and period
	GetByStrategyAndPeriod(ctx context.Context, strategyName, periodLabel string, symbol *string) (*StrategyStatistics, error)

	// GetByStrategy retrieves all statistics for a strategy
	GetByStrategy(ctx context.Context, strategyName string) ([]*StrategyStatistics, error)

	// GetByPeriod retrieves all statistics for a period
	GetByPeriod(ctx context.Context, periodLabel string) ([]*StrategyStatistics, error)

	// GetByPeriodAndStrategy retrieves statistics for a period, with optional filtering by strategy
	GetByPeriodAndStrategy(ctx context.Context, periodLabel string, strategyName *string) ([]*StrategyStatistics, error)

	// GetLatest retrieves the latest statistics for each strategy and period
	GetLatest(ctx context.Context) ([]*StrategyStatistics, error)

	// GetPreviousCalculation retrieves the statistics calculation before the given one
	// Returns the most recent calculation made before the specified calculatedAt time
	GetPreviousCalculation(ctx context.Context, strategyName, periodLabel string, symbol *string, currentCalculatedAt time.Time) (*StrategyStatistics, error)

	// GetByTimeRange retrieves statistics within a time range
	// Supports optional filtering by strategy and symbol
	GetByTimeRange(ctx context.Context, startTime, endTime time.Time, strategyName, symbol *string) ([]*StrategyStatistics, error)

	// DeleteOlderThan deletes statistics older than the specified time
	DeleteOlderThan(ctx context.Context, before time.Time) error
}
