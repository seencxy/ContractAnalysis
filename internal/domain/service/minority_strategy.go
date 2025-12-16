package service

import (
	"context"
	"fmt"

	"ContractAnalysis/internal/domain/entity"

	"github.com/shopspring/decimal"
)

// MinorityStrategyConfig represents the configuration for minority strategy
type MinorityStrategyConfig struct {
	BaseConfig                      StrategyConfig
	MinRatioDifference              float64 // Minimum ratio difference (e.g., 75 means 75:25 or more extreme)
	GenerateLongWhenShortRatioAbove float64 // Generate LONG signal when short ratio is above this
	GenerateShortWhenLongRatioAbove float64 // Generate SHORT signal when long ratio is above this
}

// MinorityStrategy implements the minority follower strategy
// Follows the minority: if 80% are short, go long
type MinorityStrategy struct {
	*BaseStrategy
	config MinorityStrategyConfig
}

// NewMinorityStrategy creates a new minority strategy
func NewMinorityStrategy(config MinorityStrategyConfig) *MinorityStrategy {
	return &MinorityStrategy{
		BaseStrategy: NewBaseStrategy(config.BaseConfig),
		config:       config,
	}
}

// Analyze analyzes market data and generates signals based on minority strategy
func (s *MinorityStrategy) Analyze(ctx context.Context, recentData []*entity.MarketData) ([]*entity.Signal, error) {
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

	// Determine signal type based on dominant direction (go opposite)
	var signalType entity.SignalType
	if latestData.GetDominantDirection() == "LONG" {
		// If majority is long, we go short
		signalType = entity.SignalTypeShort
	} else {
		// If majority is short, we go long
		signalType = entity.SignalTypeLong
	}

	// Create configuration snapshot
	configSnapshot := map[string]interface{}{
		"min_ratio_difference":                 s.config.MinRatioDifference,
		"generate_long_when_short_ratio_above": s.config.GenerateLongWhenShortRatioAbove,
		"generate_short_when_long_ratio_above": s.config.GenerateShortWhenLongRatioAbove,
		"confirmation_hours":                   s.GetConfirmationHours(),
		"tracking_hours":                       s.GetTrackingHours(),
		"profit_target_pct":                    s.GetProfitTargetPct(),
		"stop_loss_pct":                        s.GetStopLossPct(),
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
func (s *MinorityStrategy) ShouldGenerateSignal(ctx context.Context, data *entity.MarketData) (bool, string, error) {
	if !s.IsEnabled() {
		return false, "", nil
	}

	// Validate data
	if err := data.Validate(); err != nil {
		return false, "", fmt.Errorf("invalid market data: %w", err)
	}

	longThreshold := decimal.NewFromFloat(s.config.GenerateShortWhenLongRatioAbove)
	shortThreshold := decimal.NewFromFloat(s.config.GenerateLongWhenShortRatioAbove)

	// Check if short ratio is extreme (generate LONG signal)
	if data.ShortAccountRatio.GreaterThanOrEqual(shortThreshold) {
		reason := fmt.Sprintf(
			"Minority Strategy: SHORT ratio is %.2f%% (threshold: %.2f%%), going LONG to follow minority. "+
				"Long/Short ratio: %.2f%%/%.2f%%. "+
				"Position ratio: %.2f%%/%.2f%%.",
			data.ShortAccountRatio.InexactFloat64(),
			shortThreshold.InexactFloat64(),
			data.LongAccountRatio.InexactFloat64(),
			data.ShortAccountRatio.InexactFloat64(),
			data.LongPositionRatio.InexactFloat64(),
			data.ShortPositionRatio.InexactFloat64(),
		)
		return true, reason, nil
	}

	// Check if long ratio is extreme (generate SHORT signal)
	if data.LongAccountRatio.GreaterThanOrEqual(longThreshold) {
		reason := fmt.Sprintf(
			"Minority Strategy: LONG ratio is %.2f%% (threshold: %.2f%%), going SHORT to follow minority. "+
				"Long/Short ratio: %.2f%%/%.2f%%. "+
				"Position ratio: %.2f%%/%.2f%%.",
			data.LongAccountRatio.InexactFloat64(),
			longThreshold.InexactFloat64(),
			data.LongAccountRatio.InexactFloat64(),
			data.ShortAccountRatio.InexactFloat64(),
			data.LongPositionRatio.InexactFloat64(),
			data.ShortPositionRatio.InexactFloat64(),
		)
		return true, reason, nil
	}

	// Conditions not met
	return false, "", nil
}

// ValidateConfirmation checks if a signal still meets the strategy conditions
// This is used during the confirmation period to verify the signal is still valid
func (s *MinorityStrategy) ValidateConfirmation(ctx context.Context, signal *entity.Signal, currentData *entity.MarketData) (bool, string) {
	if !s.IsEnabled() {
		return false, "strategy is disabled"
	}

	// For LONG signals, check if SHORT ratio is still high
	if signal.Type == entity.SignalTypeLong {
		threshold := decimal.NewFromFloat(s.config.GenerateLongWhenShortRatioAbove)
		if currentData.ShortAccountRatio.LessThan(threshold) {
			return false, fmt.Sprintf("SHORT ratio dropped below threshold: %.2f%% < %.2f%%",
				currentData.ShortAccountRatio.InexactFloat64(),
				threshold.InexactFloat64())
		}
	}

	// For SHORT signals, check if LONG ratio is still high
	if signal.Type == entity.SignalTypeShort {
		threshold := decimal.NewFromFloat(s.config.GenerateShortWhenLongRatioAbove)
		if currentData.LongAccountRatio.LessThan(threshold) {
			return false, fmt.Sprintf("LONG ratio dropped below threshold: %.2f%% < %.2f%%",
				currentData.LongAccountRatio.InexactFloat64(),
				threshold.InexactFloat64())
		}
	}

	return true, "conditions still met"
}
