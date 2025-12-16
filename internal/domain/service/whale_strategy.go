package service

import (
	"context"
	"fmt"

	"ContractAnalysis/internal/domain/entity"

	"github.com/shopspring/decimal"
)

// WhaleStrategyConfig represents the configuration for whale strategy
type WhaleStrategyConfig struct {
	BaseConfig             StrategyConfig
	MinRatioDifference     float64 // Minimum account ratio difference
	WhalePositionThreshold float64 // Minimum whale position percentage
	MinDivergence          float64 // Minimum divergence between account and position ratios
}

// WhaleStrategy implements the whale position analysis strategy
// Detects divergence between account count ratio and position size ratio
// Example: 80% accounts long but 70% position size short -> retail being liquidated, follow whales (short)
type WhaleStrategy struct {
	*BaseStrategy
	config WhaleStrategyConfig
}

// NewWhaleStrategy creates a new whale strategy
func NewWhaleStrategy(config WhaleStrategyConfig) *WhaleStrategy {
	return &WhaleStrategy{
		BaseStrategy: NewBaseStrategy(config.BaseConfig),
		config:       config,
	}
}

// Analyze analyzes market data and generates signals based on whale strategy
func (s *WhaleStrategy) Analyze(ctx context.Context, recentData []*entity.MarketData) ([]*entity.Signal, error) {
	if !s.IsEnabled() {
		return nil, nil
	}

	if len(recentData) == 0 {
		return nil, nil
	}

	var signals []*entity.Signal

	// Analyze the most recent data point
	latestData := recentData[0]

	// Check if we should generate a signal
	shouldGenerate, reason, err := s.ShouldGenerateSignal(ctx, latestData)
	if err != nil {
		return nil, fmt.Errorf("failed to check signal condition: %w", err)
	}

	if !shouldGenerate {
		return nil, nil
	}

	// Determine signal type based on whale direction
	// We follow the whales (position ratio), not the accounts
	var signalType entity.SignalType
	if latestData.GetWhaleDirection() == "LONG" {
		signalType = entity.SignalTypeLong
	} else {
		signalType = entity.SignalTypeShort
	}

	// Create configuration snapshot
	configSnapshot := map[string]interface{}{
		"min_ratio_difference":     s.config.MinRatioDifference,
		"whale_position_threshold": s.config.WhalePositionThreshold,
		"min_divergence":           s.config.MinDivergence,
		"confirmation_hours":       s.GetConfirmationHours(),
		"tracking_hours":           s.GetTrackingHours(),
		"profit_target_pct":        s.GetProfitTargetPct(),
		"stop_loss_pct":            s.GetStopLossPct(),
	}

	// Create signal
	signal := entity.NewSignal(
		latestData.Symbol,
		signalType,
		s.Key(),
		latestData,
		s.GetConfirmationHours(),
		reason,
		configSnapshot,
	)

	signals = append(signals, signal)

	return signals, nil
}

// ShouldGenerateSignal checks if conditions are met to generate a signal
func (s *WhaleStrategy) ShouldGenerateSignal(ctx context.Context, data *entity.MarketData) (bool, string, error) {
	if !s.IsEnabled() {
		return false, "", nil
	}

	// Validate data
	if err := data.Validate(); err != nil {
		return false, "", fmt.Errorf("invalid market data: %w", err)
	}

	minRatioDiff := decimal.NewFromFloat(s.config.MinRatioDifference)
	whaleThreshold := decimal.NewFromFloat(s.config.WhalePositionThreshold)
	minDivergence := decimal.NewFromFloat(s.config.MinDivergence)

	// Check if there's divergence between account ratio and position ratio
	if !data.HasDivergence() {
		return false, "", nil
	}

	// Calculate divergence
	divergence := data.CalculateDivergence()

	// Check if divergence is significant
	if divergence.LessThan(minDivergence) {
		return false, "", nil
	}

	// Check if account ratio is extreme
	if !data.IsAccountRatioExtreme(minRatioDiff) {
		return false, "", nil
	}

	// Check if whale position meets threshold
	whaleDirection := data.GetWhaleDirection()
	var whalePositionRatio decimal.Decimal
	if whaleDirection == "LONG" {
		whalePositionRatio = data.LongPositionRatio
	} else {
		whalePositionRatio = data.ShortPositionRatio
	}

	if whalePositionRatio.LessThan(whaleThreshold) {
		return false, "", nil
	}

	// All conditions met - generate signal
	accountDirection := data.GetDominantDirection()

	reason := fmt.Sprintf(
		"Whale Strategy: Detected divergence between retail and whales. "+
			"Account ratio: %.2f%%/%.2f%% (dominant: %s). "+
			"Position ratio: %.2f%%/%.2f%% (dominant: %s). "+
			"Divergence: %.2f%% (threshold: %.2f%%). "+
			"Following whales (%s) as retail traders (%.2f%% accounts) are likely being liquidated. "+
			"Whale position: %.2f%% (threshold: %.2f%%).",
		data.LongAccountRatio.InexactFloat64(),
		data.ShortAccountRatio.InexactFloat64(),
		accountDirection,
		data.LongPositionRatio.InexactFloat64(),
		data.ShortPositionRatio.InexactFloat64(),
		whaleDirection,
		divergence.InexactFloat64(),
		minDivergence.InexactFloat64(),
		whaleDirection,
		data.GetDominantRatio().InexactFloat64(),
		whalePositionRatio.InexactFloat64(),
		whaleThreshold.InexactFloat64(),
	)

	return true, reason, nil
}

// ValidateConfirmation checks if a signal still meets the strategy conditions
func (s *WhaleStrategy) ValidateConfirmation(ctx context.Context, signal *entity.Signal, currentData *entity.MarketData) (bool, string) {
	if !s.IsEnabled() {
		return false, "strategy is disabled"
	}

	// Check if divergence still exists
	if !currentData.HasDivergence() {
		return false, "divergence no longer exists"
	}

	// Check if divergence is still significant
	minDivergence := decimal.NewFromFloat(s.config.MinDivergence)
	divergence := currentData.CalculateDivergence()
	if divergence.LessThan(minDivergence) {
		return false, fmt.Sprintf("divergence dropped below threshold: %.2f%% < %.2f%%",
			divergence.InexactFloat64(),
			minDivergence.InexactFloat64())
	}

	// Check if whale direction is still the same
	currentWhaleDirection := currentData.GetWhaleDirection()
	expectedDirection := string(signal.Type)
	if (signal.Type == entity.SignalTypeLong && currentWhaleDirection != "LONG") ||
		(signal.Type == entity.SignalTypeShort && currentWhaleDirection != "SHORT") {
		return false, fmt.Sprintf("whale direction changed from %s", expectedDirection)
	}

	// Check if whale position still meets threshold
	whaleThreshold := decimal.NewFromFloat(s.config.WhalePositionThreshold)
	var whalePositionRatio decimal.Decimal
	if signal.Type == entity.SignalTypeLong {
		whalePositionRatio = currentData.LongPositionRatio
	} else {
		whalePositionRatio = currentData.ShortPositionRatio
	}

	if whalePositionRatio.LessThan(whaleThreshold) {
		return false, fmt.Sprintf("whale position dropped below threshold: %.2f%% < %.2f%%",
			whalePositionRatio.InexactFloat64(),
			whaleThreshold.InexactFloat64())
	}

	return true, "conditions still met"
}
