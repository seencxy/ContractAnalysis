package usecase

import (
	"context"
	"fmt"
	"time"

	"ContractAnalysis/config"
	"ContractAnalysis/internal/domain/entity"
	"ContractAnalysis/internal/domain/repository"
	"ContractAnalysis/internal/domain/service"
	"ContractAnalysis/internal/infrastructure/logger"

	"go.uber.org/zap"
)

// Analyzer orchestrates signal analysis using various strategies
type Analyzer struct {
	strategies      []service.Strategy
	marketDataRepo  *repository.MarketDataRepository
	signalRepo      *repository.SignalRepository
	tradingPairRepo repository.TradingPairRepository
	globalConfig    config.GlobalStrategy
	logger          *logger.Logger
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer(
	strategies []service.Strategy,
	marketDataRepo *repository.MarketDataRepository,
	signalRepo *repository.SignalRepository,
	tradingPairRepo repository.TradingPairRepository,
	globalConfig config.GlobalStrategy,
) *Analyzer {
	return &Analyzer{
		strategies:      strategies,
		marketDataRepo:  marketDataRepo,
		signalRepo:      signalRepo,
		tradingPairRepo: tradingPairRepo,
		globalConfig:    globalConfig,
		logger:          logger.WithComponent("analyzer"),
	}
}

// AnalyzeAll analyzes market data for all trading pairs
func (a *Analyzer) AnalyzeAll(ctx context.Context) ([]*entity.Signal, error) {
	a.logger.Info("Starting signal analysis")
	startTime := time.Now()

	// Get all active trading pairs
	pairs, err := a.tradingPairRepo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active pairs: %w", err)
	}

	a.logger.Info("Analyzing trading pairs", zap.Int("count", len(pairs)))

	var allSignals []*entity.Signal

	// Analyze each pair
	for _, pair := range pairs {
		signals, err := a.analyzeSymbol(ctx, pair.Symbol)
		if err != nil {
			a.logger.WithError(err).WithSymbol(pair.Symbol).Warn("Failed to analyze symbol")
			continue
		}

		allSignals = append(allSignals, signals...)
	}

	duration := time.Since(startTime)
	a.logger.Info("Signal analysis completed",
		zap.Int("signals_generated", len(allSignals)),
		zap.String("duration", duration.String()),
	)

	return allSignals, nil
}

// AnalyzeSymbol analyzes market data for a specific symbol
func (a *Analyzer) AnalyzeSymbol(ctx context.Context, symbol string) ([]*entity.Signal, error) {
	return a.analyzeSymbol(ctx, symbol)
}

// analyzeSymbol analyzes a symbol and generates signals
func (a *Analyzer) analyzeSymbol(ctx context.Context, symbol string) ([]*entity.Signal, error) {
	mdRepo := *a.marketDataRepo
	sigRepo := *a.signalRepo

	// Get recent market data (last 24 hours)
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)
	recentData, err := mdRepo.GetBySymbol(ctx, symbol, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get market data: %w", err)
	}

	if len(recentData) == 0 {
		a.logger.Debug("No market data available for symbol", zap.String("symbol", symbol))
		return nil, nil
	}

	// Check if symbol is in cooldown period
	if inCooldown, err := a.isInCooldown(ctx, symbol); err != nil {
		return nil, fmt.Errorf("failed to check cooldown: %w", err)
	} else if inCooldown {
		a.logger.Debug("Symbol is in cooldown period", zap.String("symbol", symbol))
		return nil, nil
	}

	// Check concurrent signal limit
	if exceeded, err := a.exceedsConcurrentLimit(ctx, symbol); err != nil {
		return nil, fmt.Errorf("failed to check concurrent limit: %w", err)
	} else if exceeded {
		a.logger.Debug("Symbol has reached concurrent signal limit", zap.String("symbol", symbol))
		return nil, nil
	}

	// Apply all enabled strategies
	var allSignals []*entity.Signal

	for _, strategy := range a.strategies {
		if !strategy.IsEnabled() {
			continue
		}

		signals, err := strategy.Analyze(ctx, recentData)
		if err != nil {
			a.logger.WithError(err).WithSymbol(symbol).WithStrategy(strategy.Name()).Warn("Strategy analysis failed")
			continue
		}

		if len(signals) > 0 {
			a.logger.Info("Strategy generated signals",
				zap.String("symbol", symbol),
				zap.String("strategy", strategy.Name()),
				zap.Int("count", len(signals)),
			)

			// Store signals
			for _, signal := range signals {
				if err := sigRepo.Create(ctx, signal); err != nil {
					a.logger.WithError(err).WithSignalID(signal.SignalID).Error("Failed to store signal")
					continue
				}

				a.logger.Info("Signal created",
					zap.String("signal_id", signal.SignalID),
					zap.String("symbol", signal.Symbol),
					zap.String("type", string(signal.Type)),
					zap.String("strategy", signal.StrategyName),
				)

				allSignals = append(allSignals, signal)
			}
		}
	}

	return allSignals, nil
}

// ValidatePendingSignals validates pending signals in confirmation period
func (a *Analyzer) ValidatePendingSignals(ctx context.Context) error {
	a.logger.Info("Validating pending signals")

	sigRepo := *a.signalRepo
	mdRepo := *a.marketDataRepo

	// Get all pending signals
	pendingSignals, err := sigRepo.GetPendingSignals(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending signals: %w", err)
	}

	if len(pendingSignals) == 0 {
		a.logger.Debug("No pending signals to validate")
		return nil
	}

	a.logger.Info("Found pending signals", zap.Int("count", len(pendingSignals)))

	for _, signal := range pendingSignals {
		// Check if confirmation period has elapsed
		if !signal.ConfirmationPeriodElapsed() {
			continue
		}

		// Get latest market data
		latestData, err := mdRepo.GetLatestBySymbol(ctx, signal.Symbol)
		if err != nil {
			a.logger.WithError(err).WithSignalID(signal.SignalID).Warn("Failed to get latest market data")
			continue
		}

		if latestData == nil {
			a.logger.WithSignalID(signal.SignalID).Warn("No market data available for signal validation")
			continue
		}

		// Find the strategy that generated this signal
		var strategy service.Strategy
		for _, s := range a.strategies {
			if s.Name() == signal.StrategyName {
				strategy = s
				break
			}
		}

		if strategy == nil {
			a.logger.WithSignalID(signal.SignalID).WithStrategy(signal.StrategyName).Warn("Strategy not found for signal")
			continue
		}

		// Validate confirmation (strategy-specific validation)
		// For now, we'll just confirm the signal if it passed the confirmation period
		// In a real implementation, you'd call strategy.ValidateConfirmation()
		if err := signal.Confirm(); err != nil {
			a.logger.WithError(err).WithSignalID(signal.SignalID).Warn("Failed to confirm signal")
			continue
		}

		// Update signal
		if err := sigRepo.Update(ctx, signal); err != nil {
			a.logger.WithError(err).WithSignalID(signal.SignalID).Error("Failed to update signal")
			continue
		}

		a.logger.Info("Signal confirmed",
			zap.String("signal_id", signal.SignalID),
			zap.String("symbol", signal.Symbol),
		)
	}

	return nil
}

// isInCooldown checks if a symbol is in cooldown period
func (a *Analyzer) isInCooldown(ctx context.Context, symbol string) (bool, error) {
	if a.globalConfig.SignalCooldownHours == 0 {
		return false, nil
	}

	sigRepo := *a.signalRepo

	since := time.Now().Add(-time.Duration(a.globalConfig.SignalCooldownHours) * time.Hour)
	recentSignals, err := sigRepo.GetRecentSignalsBySymbol(ctx, symbol, since)
	if err != nil {
		return false, err
	}

	return len(recentSignals) > 0, nil
}

// exceedsConcurrentLimit checks if symbol has reached concurrent signal limit
func (a *Analyzer) exceedsConcurrentLimit(ctx context.Context, symbol string) (bool, error) {
	if a.globalConfig.MaxConcurrentSignalsPerPair == 0 {
		return false, nil
	}

	sigRepo := *a.signalRepo

	count, err := sigRepo.CountActiveSignalsBySymbol(ctx, symbol)
	if err != nil {
		return false, err
	}

	return count >= a.globalConfig.MaxConcurrentSignalsPerPair, nil
}

// GetAnalysisStatus returns the current analysis status
func (a *Analyzer) GetAnalysisStatus(ctx context.Context) (map[string]interface{}, error) {
	sigRepo := *a.signalRepo

	pendingSignals, err := sigRepo.GetPendingSignals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending signals: %w", err)
	}

	confirmedSignals, err := sigRepo.GetConfirmedSignals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get confirmed signals: %w", err)
	}

	trackingSignals, err := sigRepo.GetTrackingSignals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracking signals: %w", err)
	}

	enabledStrategies := 0
	strategyNames := []string{}
	for _, s := range a.strategies {
		if s.IsEnabled() {
			enabledStrategies++
			strategyNames = append(strategyNames, s.Name())
		}
	}

	status := map[string]interface{}{
		"enabled_strategies": enabledStrategies,
		"strategy_names":     strategyNames,
		"pending_signals":    len(pendingSignals),
		"confirmed_signals":  len(confirmedSignals),
		"tracking_signals":   len(trackingSignals),
	}

	return status, nil
}
