package usecase

import (
	"context"
	"fmt"
	"time"

	"ContractAnalysis/config"
	"ContractAnalysis/internal/domain/entity"
	"ContractAnalysis/internal/domain/repository"
	"ContractAnalysis/internal/infrastructure/logger"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// StatisticsCalculator calculates and updates strategy statistics
type StatisticsCalculator struct {
	signalRepo     *repository.SignalRepository
	statisticsRepo repository.StatisticsRepository
	config         config.StatisticsConfig
	logger         *logger.Logger
}

// NewStatisticsCalculator creates a new statistics calculator
func NewStatisticsCalculator(
	signalRepo *repository.SignalRepository,
	statisticsRepo repository.StatisticsRepository,
	cfg config.StatisticsConfig,
) *StatisticsCalculator {
	return &StatisticsCalculator{
		signalRepo:     signalRepo,
		statisticsRepo: statisticsRepo,
		config:         cfg,
		logger:         logger.WithComponent("statistics"),
	}
}

// CalculateAll calculates statistics for all strategies and periods
func (s *StatisticsCalculator) CalculateAll(ctx context.Context) error {
	s.logger.Info("Starting statistics calculation")
	startTime := time.Now()

	sigRepo := *s.signalRepo

	// Get all signals
	allSignals, err := sigRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all signals: %w", err)
	}

	if len(allSignals) == 0 {
		s.logger.Info("No signals found, skipping statistics calculation")
		return nil
	}

	s.logger.Info("Calculating statistics", zap.Int("total_signals", len(allSignals)))

	// Group signals by strategy
	signalsByStrategy := s.groupSignalsByStrategy(allSignals)

	calculated := 0
	failed := 0

	// Calculate statistics for each strategy
	for strategyName, signals := range signalsByStrategy {
		s.logger.Info("Calculating statistics for strategy",
			zap.String("strategy", strategyName),
			zap.Int("signals", len(signals)),
		)

		// Calculate for each configured period
		for _, period := range s.config.Periods {
			// Overall statistics (all symbols)
			if err := s.calculateForPeriod(ctx, strategyName, nil, signals, period); err != nil {
				s.logger.WithError(err).Error("Failed to calculate overall statistics",
					zap.String("strategy", strategyName),
					zap.String("period", period),
				)
				failed++
				continue
			}
			calculated++

			// Per-symbol statistics
			signalsBySymbol := s.groupSignalsBySymbol(signals)
			for symbol, symbolSignals := range signalsBySymbol {
				symbolCopy := symbol
				if err := s.calculateForPeriod(ctx, strategyName, &symbolCopy, symbolSignals, period); err != nil {
					s.logger.WithError(err).Warn("Failed to calculate symbol statistics",
						zap.String("strategy", strategyName),
						zap.String("symbol", symbol),
						zap.String("period", period),
					)
					failed++
					continue
				}
				calculated++
			}
		}
	}

	duration := time.Since(startTime)
	s.logger.Info("Statistics calculation completed",
		zap.Int("calculated", calculated),
		zap.Int("failed", failed),
		zap.String("duration", duration.String()),
	)

	return nil
}

// calculateForPeriod calculates statistics for a specific period
func (s *StatisticsCalculator) calculateForPeriod(
	ctx context.Context,
	strategyName string,
	symbol *string,
	signals []*entity.Signal,
	periodLabel string,
) error {
	now := time.Now()
	periodStart, periodEnd := s.getPeriodRange(now, periodLabel)

	// Filter signals by period
	periodSignals := s.filterSignalsByPeriod(signals, periodStart, periodEnd)

	if len(periodSignals) == 0 {
		// No signals in this period, skip
		return nil
	}

	// Calculate metrics
	stats := &repository.StrategyStatistics{
		StrategyName: strategyName,
		Symbol:       symbol,
		PeriodStart:  periodStart,
		PeriodEnd:    periodEnd,
		PeriodLabel:  periodLabel,
		CalculatedAt: now,
	}

	// Count signals by status
	for _, signal := range periodSignals {
		stats.TotalSignals++

		switch signal.Status {
		case entity.SignalStatusConfirmed, entity.SignalStatusTracking:
			stats.ConfirmedSignals++
		case entity.SignalStatusInvalidated:
			stats.InvalidatedSignals++
		}
	}

	// Get closed signals with outcomes
	closedSignals := s.filterClosedSignals(periodSignals)

	if len(closedSignals) > 0 {
		s.calculateOutcomeMetrics(ctx, stats, closedSignals)

		// Calculate kline metrics for closed signals
		if err := s.calculateKlineMetrics(ctx, stats, closedSignals); err != nil {
			s.logger.WithError(err).Warn("Failed to calculate kline metrics",
				zap.String("strategy", strategyName),
				zap.String("period", periodLabel),
			)
			// Don't fail the entire operation if kline metrics fail
		}
	}

	// Save statistics
	if err := s.statisticsRepo.CreateOrUpdate(ctx, stats); err != nil {
		return fmt.Errorf("failed to save statistics: %w", err)
	}

	return nil
}

// calculateOutcomeMetrics calculates performance metrics from closed signals
func (s *StatisticsCalculator) calculateOutcomeMetrics(ctx context.Context, stats *repository.StrategyStatistics, signals []*entity.Signal) {
	var totalProfit decimal.Decimal
	var totalLoss decimal.Decimal
	var totalHoldingHours decimal.Decimal
	var best *decimal.Decimal
	var worst *decimal.Decimal

	// Extract signal IDs for bulk fetching
	signalIDs := make([]string, len(signals))
	for i, signal := range signals {
		signalIDs[i] = signal.SignalID
	}

	// Fetch all outcomes in bulk
	sigRepo := *s.signalRepo
	outcomeMap, err := sigRepo.GetOutcomesBySignalIDs(ctx, signalIDs)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to fetch signal outcomes")
		outcomeMap = make(map[string]*entity.SignalOutcome)
	}

	// Process each signal with its outcome
	for _, signal := range signals {
		outcome, hasOutcome := outcomeMap[signal.SignalID]

		if !hasOutcome {
			// Signal is closed but has no outcome record (edge case)
			stats.NeutralSignals++
			s.logger.Warn("Closed signal missing outcome",
				zap.String("signal_id", signal.SignalID))
			continue
		}

		// Calculate holding hours from signal generation to close
		holdingHours := outcome.ClosedAt.Sub(signal.GeneratedAt).Hours()
		totalHoldingHours = totalHoldingHours.Add(decimal.NewFromFloat(holdingHours))

		// Classify based on actual outcome
		switch outcome.Outcome {
		case string(entity.OutcomeProfit):
			stats.ProfitableSignals++
			totalProfit = totalProfit.Add(outcome.FinalPriceChangePct)

			// Track best signal
			if best == nil || outcome.FinalPriceChangePct.GreaterThan(*best) {
				temp := outcome.FinalPriceChangePct
				best = &temp
			}

		case string(entity.OutcomeLoss):
			stats.LosingSignals++
			totalLoss = totalLoss.Add(outcome.FinalPriceChangePct.Abs())

			// Track worst signal
			if worst == nil || outcome.FinalPriceChangePct.LessThan(*worst) {
				temp := outcome.FinalPriceChangePct
				worst = &temp
			}

		default: // NEUTRAL or TIMEOUT
			stats.NeutralSignals++
		}
	}

	// Log outcome metrics for debugging
	s.logger.Debug("Outcome metrics calculated",
		zap.Int("total_closed", len(signals)),
		zap.Int("outcomes_found", len(outcomeMap)),
		zap.Int("profitable", stats.ProfitableSignals),
		zap.Int("losing", stats.LosingSignals),
		zap.Int("neutral", stats.NeutralSignals),
	)

	// Calculate averages
	if stats.ProfitableSignals > 0 {
		avgProfit := totalProfit.Div(decimal.NewFromInt(int64(stats.ProfitableSignals)))
		stats.AvgProfitPct = &avgProfit
	}

	if stats.LosingSignals > 0 {
		avgLoss := totalLoss.Div(decimal.NewFromInt(int64(stats.LosingSignals)))
		stats.AvgLossPct = &avgLoss
	}

	totalClosed := stats.ProfitableSignals + stats.LosingSignals + stats.NeutralSignals
	if totalClosed > 0 {
		avgHours := totalHoldingHours.Div(decimal.NewFromInt(int64(totalClosed)))
		stats.AvgHoldingHours = &avgHours

		// Win rate
		winRate := decimal.NewFromInt(int64(stats.ProfitableSignals)).
			Div(decimal.NewFromInt(int64(totalClosed))).
			Mul(decimal.NewFromInt(100))
		stats.WinRate = &winRate
	}

	stats.BestSignalPct = best
	stats.WorstSignalPct = worst

	// Profit factor
	if !totalLoss.IsZero() {
		profitFactor := totalProfit.Div(totalLoss)
		stats.ProfitFactor = &profitFactor
	}
}

// groupSignalsByStrategy groups signals by strategy name
func (s *StatisticsCalculator) groupSignalsByStrategy(signals []*entity.Signal) map[string][]*entity.Signal {
	groups := make(map[string][]*entity.Signal)
	for _, signal := range signals {
		groups[signal.StrategyName] = append(groups[signal.StrategyName], signal)
	}
	return groups
}

// groupSignalsBySymbol groups signals by symbol
func (s *StatisticsCalculator) groupSignalsBySymbol(signals []*entity.Signal) map[string][]*entity.Signal {
	groups := make(map[string][]*entity.Signal)
	for _, signal := range signals {
		groups[signal.Symbol] = append(groups[signal.Symbol], signal)
	}
	return groups
}

// filterSignalsByPeriod filters signals within a time period
func (s *StatisticsCalculator) filterSignalsByPeriod(signals []*entity.Signal, start, end time.Time) []*entity.Signal {
	var filtered []*entity.Signal
	for _, signal := range signals {
		if signal.GeneratedAt.After(start) && signal.GeneratedAt.Before(end) {
			filtered = append(filtered, signal)
		}
	}
	return filtered
}

// filterClosedSignals filters only closed signals
func (s *StatisticsCalculator) filterClosedSignals(signals []*entity.Signal) []*entity.Signal {
	var closed []*entity.Signal
	for _, signal := range signals {
		if signal.Status == entity.SignalStatusClosed {
			closed = append(closed, signal)
		}
	}
	return closed
}

// getPeriodRange returns the start and end time for a period label
func (s *StatisticsCalculator) getPeriodRange(now time.Time, periodLabel string) (time.Time, time.Time) {
	switch periodLabel {
	case "24h":
		return now.Add(-24 * time.Hour), now
	case "7d":
		return now.Add(-7 * 24 * time.Hour), now
	case "30d":
		return now.Add(-30 * 24 * time.Hour), now
	case "all":
		// Use a very old date for "all"
		return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), now
	default:
		// Default to 24h
		return now.Add(-24 * time.Hour), now
	}
}

// calculateKlineMetrics calculates kline-based win rate and performance metrics
func (s *StatisticsCalculator) calculateKlineMetrics(ctx context.Context, stats *repository.StrategyStatistics, signals []*entity.Signal) error {
	sigRepo := *s.signalRepo

	var totalKlineHours int
	var profitableHoursHigh int
	var profitableHoursClose int

	var sumHourlyReturn decimal.Decimal
	var maxHourlyReturn *decimal.Decimal
	var minHourlyReturn *decimal.Decimal

	var sumMaxProfit decimal.Decimal
	var sumMaxLoss decimal.Decimal

	// Process each closed signal
	for _, signal := range signals {
		if signal.Status != entity.SignalStatusClosed {
			continue
		}

		// Get all kline tracking records for this signal
		klines, err := sigRepo.GetKlineTrackingBySignal(ctx, signal.SignalID)
		if err != nil {
			s.logger.WithError(err).Warn("Failed to get kline tracking for signal",
				zap.String("signal_id", signal.SignalID),
			)
			continue
		}

		// Process each kline tracking record
		for _, kline := range klines {
			totalKlineHours++

			// Count profitable hours at high price
			if kline.IsProfitableAtHigh {
				profitableHoursHigh++
			}

			// Count profitable hours at close price
			if kline.IsProfitableAtClose {
				profitableHoursClose++
			}

			// Accumulate hourly returns
			sumHourlyReturn = sumHourlyReturn.Add(kline.HourlyReturnPct)

			// Track max/min hourly returns
			if maxHourlyReturn == nil || kline.HourlyReturnPct.GreaterThan(*maxHourlyReturn) {
				temp := kline.HourlyReturnPct
				maxHourlyReturn = &temp
			}

			if minHourlyReturn == nil || kline.HourlyReturnPct.LessThan(*minHourlyReturn) {
				temp := kline.HourlyReturnPct
				minHourlyReturn = &temp
			}

			// Accumulate theoretical max profit/loss
			sumMaxProfit = sumMaxProfit.Add(kline.MaxPotentialProfitPct)
			sumMaxLoss = sumMaxLoss.Add(kline.MaxPotentialLossPct)
		}
	}

	// Calculate and set kline statistics
	stats.TotalKlineHours = totalKlineHours
	stats.ProfitableKlineHoursHigh = profitableHoursHigh
	stats.ProfitableKlineHoursClose = profitableHoursClose

	if totalKlineHours > 0 {
		// Theoretical win rate (based on high price)
		theoreticalWinRate := decimal.NewFromInt(int64(profitableHoursHigh)).
			Div(decimal.NewFromInt(int64(totalKlineHours))).
			Mul(decimal.NewFromInt(100))
		stats.KlineTheoreticalWinRate = &theoreticalWinRate

		// Close win rate (based on close price)
		closeWinRate := decimal.NewFromInt(int64(profitableHoursClose)).
			Div(decimal.NewFromInt(int64(totalKlineHours))).
			Mul(decimal.NewFromInt(100))
		stats.KlineCloseWinRate = &closeWinRate

		// Average hourly return
		avgHourlyReturn := sumHourlyReturn.Div(decimal.NewFromInt(int64(totalKlineHours)))
		stats.AvgHourlyReturnPct = &avgHourlyReturn

		// Max/Min hourly return
		stats.MaxHourlyReturnPct = maxHourlyReturn
		stats.MinHourlyReturnPct = minHourlyReturn

		// Average maximum potential profit/loss
		avgMaxProfit := sumMaxProfit.Div(decimal.NewFromInt(int64(totalKlineHours)))
		stats.AvgMaxPotentialProfitPct = &avgMaxProfit

		avgMaxLoss := sumMaxLoss.Div(decimal.NewFromInt(int64(totalKlineHours)))
		stats.AvgMaxPotentialLossPct = &avgMaxLoss
	}

	return nil
}
