package usecase

import (
	"context"
	"fmt"
	"time"

	"ContractAnalysis/internal/domain/entity"
	"ContractAnalysis/internal/domain/repository"
	"ContractAnalysis/internal/infrastructure/binance"
	"ContractAnalysis/internal/infrastructure/logger"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// Tracker orchestrates signal tracking and outcome calculation
type Tracker struct {
	binanceClient *binance.Client
	signalRepo    *repository.SignalRepository
	logger        *logger.Logger
}

// NewTracker creates a new tracker
func NewTracker(
	binanceClient *binance.Client,
	signalRepo *repository.SignalRepository,
) *Tracker {
	return &Tracker{
		binanceClient: binanceClient,
		signalRepo:    signalRepo,
		logger:        logger.WithComponent("tracker"),
	}
}

// TrackAll tracks all active signals
func (t *Tracker) TrackAll(ctx context.Context) error {
	t.logger.Info("Starting signal tracking")
	startTime := time.Now()

	sigRepo := *t.signalRepo

	// Get all signals that need tracking (CONFIRMED and TRACKING)
	confirmedSignals, err := sigRepo.GetConfirmedSignals(ctx)
	if err != nil {
		return fmt.Errorf("failed to get confirmed signals: %w", err)
	}

	trackingSignals, err := sigRepo.GetTrackingSignals(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tracking signals: %w", err)
	}

	allSignals := append(confirmedSignals, trackingSignals...)

	if len(allSignals) == 0 {
		t.logger.Debug("No signals to track")
		return nil
	}

	t.logger.Info("Tracking signals", zap.Int("count", len(allSignals)))

	tracked := 0
	completed := 0
	failed := 0

	for _, signal := range allSignals {
		if err := t.trackSignal(ctx, signal); err != nil {
			t.logger.WithError(err).WithSignalID(signal.SignalID).Warn("Failed to track signal")
			failed++
			continue
		}

		tracked++

		// Check if signal is completed
		if signal.Status == entity.SignalStatusClosed {
			completed++
		}
	}

	duration := time.Since(startTime)
	t.logger.Info("Signal tracking completed",
		zap.Int("tracked", tracked),
		zap.Int("completed", completed),
		zap.Int("failed", failed),
		zap.String("duration", duration.String()),
	)

	return nil
}

// TrackSignal tracks a specific signal
func (t *Tracker) TrackSignal(ctx context.Context, signalID string) error {
	sigRepo := *t.signalRepo

	signal, err := sigRepo.GetByID(ctx, signalID)
	if err != nil {
		return fmt.Errorf("failed to get signal: %w", err)
	}

	if signal == nil {
		return fmt.Errorf("signal not found: %s", signalID)
	}

	return t.trackSignal(ctx, signal)
}

// trackSignal tracks a signal and updates its status
func (t *Tracker) trackSignal(ctx context.Context, signal *entity.Signal) error {
	sigRepo := *t.signalRepo

	// Get current price
	currentPrice, err := t.binanceClient.GetPrice(ctx, signal.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get current price: %w", err)
	}

	currentPriceDecimal := decimal.NewFromFloat(currentPrice)

	// Get latest tracking record
	latestTracking, err := sigRepo.GetLatestTracking(ctx, signal.SignalID)
	if err != nil {
		return fmt.Errorf("failed to get latest tracking: %w", err)
	}

	// Calculate price change
	priceChangePct := signal.CalculatePriceChange(currentPriceDecimal)

	// Create or update tracking record
	var tracking *entity.SignalTracking
	if latestTracking == nil {
		// First tracking record
		tracking = entity.NewSignalTracking(signal.SignalID, signal, currentPriceDecimal)
	} else {
		// Create new tracking record
		tracking = entity.NewSignalTracking(signal.SignalID, signal, currentPriceDecimal)

		// Update peak/trough from previous tracking
		tracking.HighestPrice = latestTracking.HighestPrice
		tracking.HighestPricePct = latestTracking.HighestPricePct
		tracking.HighestPriceAt = latestTracking.HighestPriceAt
		tracking.LowestPrice = latestTracking.LowestPrice
		tracking.LowestPricePct = latestTracking.LowestPricePct
		tracking.LowestPriceAt = latestTracking.LowestPriceAt

		// Update if new peak or trough
		tracking.UpdatePeakTrough(currentPriceDecimal, priceChangePct)
	}

	// Save tracking record
	if err := sigRepo.CreateTracking(ctx, tracking); err != nil {
		return fmt.Errorf("failed to create tracking: %w", err)
	}

	t.logger.Debug("Tracking updated",
		zap.String("signal_id", signal.SignalID),
		zap.String("price_change", priceChangePct.String()),
	)

	// Update signal status if needed
	if signal.Status == entity.SignalStatusConfirmed {
		// Start tracking
		if err := signal.StartTracking(); err != nil {
			return fmt.Errorf("failed to start tracking: %w", err)
		}
		if err := sigRepo.Update(ctx, signal); err != nil {
			return fmt.Errorf("failed to update signal: %w", err)
		}
		t.logger.Info("Signal tracking started", zap.String("signal_id", signal.SignalID))
	}

	// Check if signal should be closed
	// Get strategy config from signal
	profitTargetPct := decimal.NewFromFloat(5.0) // Default
	stopLossPct := decimal.NewFromFloat(2.0)     // Default
	trackingHours := 24                          // Default

	if signal.ConfigSnapshot != nil {
		if val, ok := signal.ConfigSnapshot["profit_target_pct"].(float64); ok {
			profitTargetPct = decimal.NewFromFloat(val)
		}
		if val, ok := signal.ConfigSnapshot["stop_loss_pct"].(float64); ok {
			stopLossPct = decimal.NewFromFloat(val)
		}
		if val, ok := signal.ConfigSnapshot["tracking_hours"].(float64); ok {
			trackingHours = int(val)
		}
	}

	shouldClose := false
	closeReason := ""

	// --- Enhanced Trade Management Logic ---

	// Check Stop Loss (Priority)
	// If dynamic SL is set, use it. Otherwise use percentage based.
	isStopLossHit := false
	if !signal.StopLossPrice.IsZero() {
		// For SHORT: Hit if Price >= SL
		if signal.Type == entity.SignalTypeShort && currentPriceDecimal.GreaterThanOrEqual(signal.StopLossPrice) {
			isStopLossHit = true
		}
		// For LONG: Hit if Price <= SL
		if signal.Type == entity.SignalTypeLong && currentPriceDecimal.LessThanOrEqual(signal.StopLossPrice) {
			isStopLossHit = true
		}
	} else {
		// Fallback to percentage SL
		if priceChangePct.LessThanOrEqual(stopLossPct.Neg()) {
			isStopLossHit = true
		}
	}

	if isStopLossHit {
		shouldClose = true
		closeReason = "stop loss hit"
		signal.ExitPrice = currentPriceDecimal
		signal.ExitReason = "SL"
	} else {
		// Check Take Profit 1
		// TODO: Implement Partial Close logic (requires Order/Position entity)
		// For now, we just log it or maybe move SL to breakeven?

		// Check Final Take Profit (TP2 or Percentage)
		isTPHit := false
		if !signal.TargetPrice2.IsZero() {
			// For SHORT: Hit if Price <= TP2
			if signal.Type == entity.SignalTypeShort && currentPriceDecimal.LessThanOrEqual(signal.TargetPrice2) {
				isTPHit = true
			}
			// For LONG: Hit if Price >= TP2
			if signal.Type == entity.SignalTypeLong && currentPriceDecimal.GreaterThanOrEqual(signal.TargetPrice2) {
				isTPHit = true
			}
		} else {
			// Fallback to percentage TP
			if priceChangePct.GreaterThanOrEqual(profitTargetPct) {
				isTPHit = true
			}
		}

		if isTPHit {
			shouldClose = true
			closeReason = "profit target reached"
			signal.ExitPrice = currentPriceDecimal
			signal.ExitReason = "TP"
		}
	}

	// Check tracking time limit
	if !shouldClose && signal.HoursElapsed() >= float64(trackingHours) {
		shouldClose = true
		closeReason = "tracking period elapsed"
		signal.ExitPrice = currentPriceDecimal
		signal.ExitReason = "Time"
	}

	if shouldClose && signal.Status == entity.SignalStatusTracking {
		// Close signal and calculate outcome
		if err := signal.Close(); err != nil {
			return fmt.Errorf("failed to close signal: %w", err)
		}

		// Create outcome
		outcome := entity.NewSignalOutcome(signal.SignalID, signal, tracking, profitTargetPct, stopLossPct)
		if err := sigRepo.CreateOutcome(ctx, outcome); err != nil {
			return fmt.Errorf("failed to create outcome: %w", err)
		}

		// Update signal
		if err := sigRepo.Update(ctx, signal); err != nil {
			return fmt.Errorf("failed to update signal: %w", err)
		}

		t.logger.Info("Signal closed",
			zap.String("signal_id", signal.SignalID),
			zap.String("reason", closeReason),
			zap.String("outcome", outcome.Outcome),
			zap.String("final_change", outcome.FinalPriceChangePct.String()),
		)
	}

	return nil
}

// GetTrackingStatus returns the current tracking status
func (t *Tracker) GetTrackingStatus(ctx context.Context) (map[string]interface{}, error) {
	sigRepo := *t.signalRepo

	trackingSignals, err := sigRepo.GetTrackingSignals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracking signals: %w", err)
	}

	confirmedSignals, err := sigRepo.GetConfirmedSignals(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get confirmed signals: %w", err)
	}

	// Get recent outcomes (last 24 hours)
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)
	recentOutcomes, err := sigRepo.GetOutcomesByTimeRange(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent outcomes: %w", err)
	}

	// Calculate outcome statistics
	profitable := 0
	losing := 0
	for _, outcome := range recentOutcomes {
		if outcome.IsProfit() {
			profitable++
		} else if outcome.IsLoss() {
			losing++
		}
	}

	status := map[string]interface{}{
		"tracking_signals":  len(trackingSignals),
		"confirmed_signals": len(confirmedSignals),
		"recent_outcomes":   len(recentOutcomes),
		"recent_profitable": profitable,
		"recent_losing":     losing,
	}

	return status, nil
}

// TrackAllKlines tracks all active signals using kline data (runs hourly)
func (t *Tracker) TrackAllKlines(ctx context.Context) error {
	t.logger.Info("Starting kline tracking")
	startTime := time.Now()

	sigRepo := *t.signalRepo

	// Get all signals that need kline tracking (CONFIRMED and TRACKING)
	confirmedSignals, err := sigRepo.GetConfirmedSignals(ctx)
	if err != nil {
		return fmt.Errorf("failed to get confirmed signals: %w", err)
	}

	trackingSignals, err := sigRepo.GetTrackingSignals(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tracking signals: %w", err)
	}

	allSignals := append(confirmedSignals, trackingSignals...)

	if len(allSignals) == 0 {
		t.logger.Debug("No signals for kline tracking")
		return nil
	}

	t.logger.Info("Tracking klines for signals", zap.Int("count", len(allSignals)))

	// Group signals by symbol for batch optimization
	signalsBySymbol := make(map[string][]*entity.Signal)
	for _, signal := range allSignals {
		signalsBySymbol[signal.Symbol] = append(signalsBySymbol[signal.Symbol], signal)
	}

	tracked := 0
	failed := 0

	// Process each symbol's signals
	for symbol, signals := range signalsBySymbol {
		if err := t.trackSymbolKlines(ctx, symbol, signals); err != nil {
			t.logger.WithError(err).WithSymbol(symbol).Warn("Failed to track klines for symbol")
			failed += len(signals)
			continue
		}
		tracked += len(signals)

		// Avoid API rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	duration := time.Since(startTime)
	t.logger.Info("Kline tracking completed",
		zap.Int("tracked", tracked),
		zap.Int("failed", failed),
		zap.String("duration", duration.String()),
	)

	return nil
}

// trackSymbolKlines tracks klines for all signals of a specific symbol
func (t *Tracker) trackSymbolKlines(ctx context.Context, symbol string, signals []*entity.Signal) error {
	sigRepo := *t.signalRepo

	// Determine the earliest start time among all signals
	var earliestStart time.Time
	for i, signal := range signals {
		// Get latest kline tracking for this signal
		latestKline, err := sigRepo.GetLatestKlineTracking(ctx, signal.SignalID)
		if err != nil {
			return fmt.Errorf("failed to get latest kline tracking: %w", err)
		}

		var startTime time.Time
		if latestKline != nil {
			// Start from next hour after last tracked kline
			startTime = latestKline.KlineCloseTime.Add(1 * time.Second)
		} else {
			// Start from signal generation time (truncate to hour)
			startTime = signal.GeneratedAt.Truncate(time.Hour)
		}

		if i == 0 || startTime.Before(earliestStart) {
			earliestStart = startTime
		}
	}

	// Get current hour (don't fetch incomplete kline)
	now := time.Now()
	currentHour := now.Truncate(time.Hour)

	// Skip if no complete klines available
	if earliestStart.After(currentHour) || earliestStart.Equal(currentHour) {
		t.logger.Debug("No new complete klines to track",
			zap.String("symbol", symbol),
		)
		return nil
	}

	// Fetch klines for this symbol
	klines, err := t.binanceClient.GetKlinesSince(ctx, symbol, "1h", earliestStart)
	if err != nil {
		return fmt.Errorf("failed to get klines: %w", err)
	}

	// Filter out incomplete klines (current hour)
	var completedKlines []*entity.Kline
	for _, kline := range klines {
		if kline.CloseTime.Before(currentHour) {
			completedKlines = append(completedKlines, kline)
		}
	}

	if len(completedKlines) == 0 {
		t.logger.Debug("No completed klines available",
			zap.String("symbol", symbol),
		)
		return nil
	}

	// Process klines for each signal
	for _, signal := range signals {
		if err := t.processSignalKlines(ctx, signal, completedKlines); err != nil {
			t.logger.WithError(err).WithSignalID(signal.SignalID).Warn("Failed to process klines for signal")
			continue
		}
	}

	return nil
}

// processSignalKlines processes klines for a single signal
func (t *Tracker) processSignalKlines(ctx context.Context, signal *entity.Signal, klines []*entity.Kline) error {
	sigRepo := *t.signalRepo

	// Get latest kline tracking to avoid duplicates
	latestKline, err := sigRepo.GetLatestKlineTracking(ctx, signal.SignalID)
	if err != nil {
		return fmt.Errorf("failed to get latest kline tracking: %w", err)
	}

	var lastTrackedTime time.Time
	if latestKline != nil {
		lastTrackedTime = latestKline.KlineCloseTime
	}

	// Create kline tracking records for new klines only
	for _, kline := range klines {
		// Skip if kline is before signal generation or already tracked
		if kline.OpenTime.Before(signal.GeneratedAt) || kline.CloseTime.Before(lastTrackedTime) || kline.CloseTime.Equal(lastTrackedTime) {
			continue
		}

		// Create kline tracking record
		tracking := entity.NewSignalKlineTracking(signal.SignalID, signal, kline)

		if err := sigRepo.CreateKlineTracking(ctx, tracking); err != nil {
			return fmt.Errorf("failed to create kline tracking: %w", err)
		}

		t.logger.Debug("Kline tracking created",
			zap.String("signal_id", signal.SignalID),
			zap.Time("kline_time", kline.OpenTime),
			zap.String("close_change", tracking.CloseChangePct.String()),
		)
	}

	return nil
}
