package usecase

import (
	"context"
	"fmt"
	"math"

	"ContractAnalysis/config"
	"ContractAnalysis/internal/domain/repository"
	"ContractAnalysis/internal/infrastructure/logger"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// StatisticsMonitor monitors statistics changes and alerts on significant variations
type StatisticsMonitor struct {
	statisticsRepo repository.StatisticsRepository
	config         config.StatisticsMonitoringConfig
	logger         *logger.Logger
}

// MetricChange represents a detected change in a metric
type MetricChange struct {
	MetricName     string
	PreviousValue  string
	CurrentValue   string
	Change         float64
	ChangeType     string // "percentage" or "percentage_points"
	IsSignificant  bool
}

// NewStatisticsMonitor creates a new statistics monitor
func NewStatisticsMonitor(
	statisticsRepo repository.StatisticsRepository,
	config config.StatisticsMonitoringConfig,
) *StatisticsMonitor {
	return &StatisticsMonitor{
		statisticsRepo: statisticsRepo,
		config:         config,
		logger:         logger.WithComponent("statistics_monitor"),
	}
}

// MonitorAllStatistics monitors all latest statistics for changes
func (m *StatisticsMonitor) MonitorAllStatistics(ctx context.Context) error {
	if !m.config.Enabled {
		m.logger.Debug("Statistics monitoring is disabled")
		return nil
	}

	// Get all latest statistics
	allStats, err := m.statisticsRepo.GetLatest(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest statistics: %w", err)
	}

	if len(allStats) == 0 {
		m.logger.Debug("No statistics found to monitor")
		return nil
	}

	monitored := 0
	warnings := 0

	for _, current := range allStats {
		changes, err := m.MonitorChanges(ctx, current)
		if err != nil {
			m.logger.WithError(err).Warn("Failed to monitor statistics",
				zap.String("strategy", current.StrategyName),
				zap.String("period", current.PeriodLabel))
			continue
		}

		monitored++
		if len(changes) > 0 {
			warnings++
		}
	}

	m.logger.Info("Statistics monitoring completed",
		zap.Int("monitored", monitored),
		zap.Int("warnings", warnings))

	return nil
}

// MonitorChanges monitors a single statistics record for changes
func (m *StatisticsMonitor) MonitorChanges(
	ctx context.Context,
	current *repository.StrategyStatistics,
) ([]MetricChange, error) {
	// Get previous calculation
	previous, err := m.statisticsRepo.GetPreviousCalculation(
		ctx,
		current.StrategyName,
		current.PeriodLabel,
		current.Symbol,
		current.CalculatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get previous calculation: %w", err)
	}

	if previous == nil {
		// First calculation, nothing to compare
		m.logger.Debug("No previous calculation found, skipping monitoring",
			zap.String("strategy", current.StrategyName),
			zap.String("period", current.PeriodLabel))
		return nil, nil
	}

	// Check for significant changes
	changes := m.detectSignificantChanges(current, previous)

	if len(changes) > 0 {
		m.logChanges(current, previous, changes)
	}

	return changes, nil
}

// detectSignificantChanges compares current and previous statistics
func (m *StatisticsMonitor) detectSignificantChanges(
	current, previous *repository.StrategyStatistics,
) []MetricChange {
	var changes []MetricChange

	// Check Win Rate (percentage point change)
	if current.WinRate != nil && previous.WinRate != nil {
		change := current.WinRate.Sub(*previous.WinRate).InexactFloat64()
		if math.Abs(change) >= m.config.WinRateChangeThreshold {
			changes = append(changes, MetricChange{
				MetricName:    "Win Rate",
				PreviousValue: fmt.Sprintf("%.2f%%", previous.WinRate.InexactFloat64()),
				CurrentValue:  fmt.Sprintf("%.2f%%", current.WinRate.InexactFloat64()),
				Change:        change,
				ChangeType:    "percentage_points",
				IsSignificant: true,
			})
		}
	}

	// Check Profitable Signals Ratio (percentage point change)
	if current.TotalSignals > 0 && previous.TotalSignals > 0 {
		prevRatio := float64(previous.ProfitableSignals) / float64(previous.TotalSignals) * 100
		currRatio := float64(current.ProfitableSignals) / float64(current.TotalSignals) * 100
		change := currRatio - prevRatio

		if math.Abs(change) >= m.config.ProfitRatioChangeThreshold {
			changes = append(changes, MetricChange{
				MetricName:    "Profitable Signals Ratio",
				PreviousValue: fmt.Sprintf("%.2f%%", prevRatio),
				CurrentValue:  fmt.Sprintf("%.2f%%", currRatio),
				Change:        change,
				ChangeType:    "percentage_points",
				IsSignificant: true,
			})
		}
	}

	// Check Average Profit (percentage change)
	if current.AvgProfitPct != nil && previous.AvgProfitPct != nil && !previous.AvgProfitPct.IsZero() {
		percentChange := current.AvgProfitPct.Sub(*previous.AvgProfitPct).
			Div(*previous.AvgProfitPct).
			Mul(decimal.NewFromInt(100)).
			InexactFloat64()

		if math.Abs(percentChange) >= m.config.AvgProfitChangeThreshold {
			changes = append(changes, MetricChange{
				MetricName:    "Average Profit",
				PreviousValue: fmt.Sprintf("%.2f%%", previous.AvgProfitPct.InexactFloat64()),
				CurrentValue:  fmt.Sprintf("%.2f%%", current.AvgProfitPct.InexactFloat64()),
				Change:        percentChange,
				ChangeType:    "percentage",
				IsSignificant: true,
			})
		}
	}

	// Check Average Loss (percentage change)
	if current.AvgLossPct != nil && previous.AvgLossPct != nil && !previous.AvgLossPct.IsZero() {
		percentChange := current.AvgLossPct.Sub(*previous.AvgLossPct).
			Div(*previous.AvgLossPct).
			Mul(decimal.NewFromInt(100)).
			InexactFloat64()

		if math.Abs(percentChange) >= m.config.AvgLossChangeThreshold {
			changes = append(changes, MetricChange{
				MetricName:    "Average Loss",
				PreviousValue: fmt.Sprintf("%.2f%%", previous.AvgLossPct.InexactFloat64()),
				CurrentValue:  fmt.Sprintf("%.2f%%", current.AvgLossPct.InexactFloat64()),
				Change:        percentChange,
				ChangeType:    "percentage",
				IsSignificant: true,
			})
		}
	}

	// Check Profit Factor (percentage change)
	if current.ProfitFactor != nil && previous.ProfitFactor != nil && !previous.ProfitFactor.IsZero() {
		percentChange := current.ProfitFactor.Sub(*previous.ProfitFactor).
			Div(*previous.ProfitFactor).
			Mul(decimal.NewFromInt(100)).
			InexactFloat64()

		if math.Abs(percentChange) >= m.config.ProfitFactorChangeThreshold {
			changes = append(changes, MetricChange{
				MetricName:    "Profit Factor",
				PreviousValue: fmt.Sprintf("%.2f", previous.ProfitFactor.InexactFloat64()),
				CurrentValue:  fmt.Sprintf("%.2f", current.ProfitFactor.InexactFloat64()),
				Change:        percentChange,
				ChangeType:    "percentage",
				IsSignificant: true,
			})
		}
	}

	// Check Total Signals (percentage change)
	if previous.TotalSignals > 0 {
		percentChange := float64(current.TotalSignals-previous.TotalSignals) / float64(previous.TotalSignals) * 100

		if math.Abs(percentChange) >= m.config.SignalCountChangeThreshold {
			changes = append(changes, MetricChange{
				MetricName:    "Total Signals",
				PreviousValue: fmt.Sprintf("%d", previous.TotalSignals),
				CurrentValue:  fmt.Sprintf("%d", current.TotalSignals),
				Change:        percentChange,
				ChangeType:    "percentage",
				IsSignificant: true,
			})
		}
	}

	return changes
}

// logChanges logs detected changes in a formatted message
func (m *StatisticsMonitor) logChanges(
	current, previous *repository.StrategyStatistics,
	changes []MetricChange,
) {
	symbolStr := "ALL"
	if current.Symbol != nil {
		symbolStr = *current.Symbol
	}

	message := fmt.Sprintf(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  STATISTICS CHANGE DETECTED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Strategy:  %s
Symbol:    %s
Period:    %s
Previous:  %s
Current:   %s

Significant Changes:
`, current.StrategyName, symbolStr, current.PeriodLabel,
		previous.CalculatedAt.Format("2006-01-02 15:04:05"),
		current.CalculatedAt.Format("2006-01-02 15:04:05"))

	for _, change := range changes {
		changeSymbol := "📈"
		if change.Change < 0 {
			changeSymbol = "📉"
		}

		changeTypeStr := "change"
		if change.ChangeType == "percentage_points" {
			changeTypeStr = "point change"
		}

		message += fmt.Sprintf("%s %s: %s → %s (%+.2f%% %s)\n",
			changeSymbol,
			change.MetricName,
			change.PreviousValue,
			change.CurrentValue,
			change.Change,
			changeTypeStr)
	}

	message += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

	m.logger.Warn(message)
}
