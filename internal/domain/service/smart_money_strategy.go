package service

import (
	"context"
	"fmt"

	"ContractAnalysis/internal/domain/entity"
	"ContractAnalysis/internal/domain/repository"

	"github.com/shopspring/decimal"
)

// SmartMoneyStrategyConfig represents the configuration for Smart Money strategy
type SmartMoneyStrategyConfig struct {
	BaseConfig          StrategyConfig
	MinLongAccountRatio float64 // Minimum Long Account Ratio to consider (e.g. 70%)
	LookbackPeriod      int     // Number of candles to look back for High (e.g. 24)
	KlineInterval       string  // Interval for klines (e.g. "1h", "15m")
}

// SmartMoneyStrategy implements the "Three-Step" Smart Money logic
// 1. Monitor: High Retail Long Ratio (>2.0 or 66%), Rising OI, Smart Money divergence
// 2. Trigger: Swing Failure Pattern (SFP) / Liquidity Grab at previous high
// 3. Exit: Managed by BaseStrategy (Stop Loss above fake-out high, Profit Target at low)
type SmartMoneyStrategy struct {
	*BaseStrategy
	config          SmartMoneyStrategyConfig
	klineRepo       repository.KlineRepository
	patternAnalyzer *PatternAnalyzer
}

// NewSmartMoneyStrategy creates a new Smart Money strategy
func NewSmartMoneyStrategy(config SmartMoneyStrategyConfig, klineRepo repository.KlineRepository) *SmartMoneyStrategy {
	return &SmartMoneyStrategy{
		BaseStrategy:    NewBaseStrategy(config.BaseConfig),
		config:          config,
		klineRepo:       klineRepo,
		patternAnalyzer: NewPatternAnalyzer(),
	}
}

// Analyze analyzes market data and generates signals
func (s *SmartMoneyStrategy) Analyze(ctx context.Context, recentData []*entity.MarketData) ([]*entity.Signal, error) {
	if !s.IsEnabled() {
		return nil, nil
	}

	if len(recentData) == 0 {
		return nil, nil
	}

	// Use the most recent data point
	latestData := recentData[0]

	// Check if we should generate a signal
	shouldGenerate, _, err := s.ShouldGenerateSignal(ctx, latestData)
	if err != nil {
		return nil, fmt.Errorf("failed to check signal condition: %w", err)
	}

	if !shouldGenerate {
		return nil, nil
	}

	// Smart Money Logic implies we are Shorting the liquidity grab
	signalType := entity.SignalTypeShort

	setup, err := s.detectSFPSetup(ctx, latestData)
	if err != nil {
		return nil, fmt.Errorf("failed to detect setup: %w", err)
	}

	if setup == nil {
		return nil, nil
	}

	// Create configuration snapshot
	configSnapshot := map[string]interface{}{
		"min_long_account_ratio": s.config.MinLongAccountRatio,
		"lookback_period":        s.config.LookbackPeriod,
		"kline_interval":         s.config.KlineInterval,
		"confirmation_hours":     s.GetConfirmationHours(),
		"tracking_hours":         s.GetTrackingHours(),
		"profit_target_pct":      s.GetProfitTargetPct(),
		"stop_loss_pct":          s.GetStopLossPct(),
		"setup_type":             "SFP_SHORT",
	}

	// Create signal
	signal := entity.NewSignal(
		latestData.Symbol,
		signalType,
		s.Key(),
		latestData,
		s.GetConfirmationHours(),
		setup.Reason,
		configSnapshot,
	)

	// Set Trade Levels
	signal.SetTradeLevels(setup.StopLoss, setup.TakeProfit1, setup.TakeProfit2)

	// Enable trailing stop if configured
	trailingStopCfg := s.GetTrailingStopConfig()
	if trailingStopCfg.Enabled {
		signal.TrailingStopEnabled = true
		signal.TrailingStopActivationPct = decimal.NewFromFloat(trailingStopCfg.ActivationPct)
		signal.TrailingStopDistancePct = decimal.NewFromFloat(trailingStopCfg.TrailDistancePct)
	}

	return []*entity.Signal{signal}, nil
}

// TradeSetup holds calculated trade parameters
type TradeSetup struct {
	Reason      string
	StopLoss    decimal.Decimal
	TakeProfit1 decimal.Decimal
	TakeProfit2 decimal.Decimal
}

// detectSFPSetup performs the full SFP detection and calculates trade levels
func (s *SmartMoneyStrategy) detectSFPSetup(ctx context.Context, data *entity.MarketData) (*TradeSetup, error) {
	// Fetch Klines
	klines, err := s.klineRepo.GetKlines(ctx, data.Symbol, s.config.KlineInterval, s.config.LookbackPeriod+2)
	if err != nil {
		return nil, err
	}
	if len(klines) < 3 {
		return nil, nil
	}

	triggerCandle := klines[len(klines)-2]
	prevCandle := klines[len(klines)-3]

	// Find Swing High
	startIdx := len(klines) - 2 - s.config.LookbackPeriod
	if startIdx < 0 {
		startIdx = 0
	}
	endIdx := len(klines) - 3

	highestHigh := decimal.Zero
	lowestLow := decimal.NewFromFloat(1000000000.0) // Arbitrary large number

	for i := startIdx; i <= endIdx; i++ {
		if klines[i].High.GreaterThan(highestHigh) {
			highestHigh = klines[i].High
		}
		if klines[i].Low.LessThan(lowestLow) {
			lowestLow = klines[i].Low
		}
	}

	// Pattern Check (Confluence)
	isSFP := s.patternAnalyzer.IsSwingFailurePattern(triggerCandle, highestHigh)
	isShootingStar := s.patternAnalyzer.IsShootingStar(triggerCandle)
	isBearishEngulfing := s.patternAnalyzer.IsBearishEngulfing(triggerCandle, prevCandle)

	if isSFP || isShootingStar || isBearishEngulfing {
		// Calculate SL: High of the trigger candle + buffer (0.1%)
		stopLoss := triggerCandle.High.Mul(decimal.NewFromFloat(1.001))

		// If Bearish Engulfing, SL can be the high of the engulfing candle or the previous one (whichever is higher)
		if isBearishEngulfing && prevCandle.High.GreaterThan(triggerCandle.High) {
			stopLoss = prevCandle.High.Mul(decimal.NewFromFloat(1.001))
		}

		// Calculate TP1: Lowest Low of lookback period
		takeProfit1 := lowestLow

		// Calculate TP2: Fixed Risk:Reward of 1:3
		entryPrice := triggerCandle.Close
		risk := stopLoss.Sub(entryPrice)

		// Ensure risk is positive (SL > Entry for Short)
		if risk.LessThanOrEqual(decimal.Zero) {
			// Fallback if SL is too close or invalid logic
			risk = entryPrice.Mul(decimal.NewFromFloat(0.01))
		}

		takeProfit2 := entryPrice.Sub(risk.Mul(decimal.NewFromFloat(3.0)))

		patternName := ""
		if isSFP {
			patternName += "SFP "
		}
		if isShootingStar {
			patternName += "ShootingStar "
		}
		if isBearishEngulfing {
			patternName += "BearishEngulfing "
		}

		reason := fmt.Sprintf("Smart Money Confluence: %sdetected. SL at %.2f, TP1 at %.2f (Low)", patternName, stopLoss.InexactFloat64(), takeProfit1.InexactFloat64())

		return &TradeSetup{
			Reason:      reason,
			StopLoss:    stopLoss,
			TakeProfit1: takeProfit1,
			TakeProfit2: takeProfit2,
		}, nil
	}

	return nil, nil
}

// ShouldGenerateSignal checks if conditions are met to generate a signal
func (s *SmartMoneyStrategy) ShouldGenerateSignal(ctx context.Context, data *entity.MarketData) (bool, string, error) {
	if !s.IsEnabled() {
		return false, "", nil
	}

	if err := data.Validate(); err != nil {
		return false, "", fmt.Errorf("invalid market data: %w", err)
	}

	// 1.1 Long Account Ratio Check
	minLongRatio := decimal.NewFromFloat(s.config.MinLongAccountRatio)
	if data.LongAccountRatio.LessThan(minLongRatio) {
		return false, "", nil
	}

	// 1.2 Divergence Check
	if data.LongPositionRatio.GreaterThanOrEqual(data.LongAccountRatio) {
		return false, "", nil
	}

	// Funding Rate Check (Positive) - Crowd is Long paying Short
	if data.FundingRate.LessThanOrEqual(decimal.Zero) {
		return false, "", nil
	}

	// Call Setup Detector
	setup, err := s.detectSFPSetup(ctx, data)
	if err != nil {
		return false, "", err
	}

	if setup != nil {
		return true, setup.Reason, nil
	}

	return false, "", nil
}
