package entity

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// SignalTracking represents tracking data for a signal
type SignalTracking struct {
	ID       int64
	SignalID string

	// Tracking details
	TrackedAt    time.Time
	HoursElapsed decimal.Decimal

	CurrentPrice    decimal.Decimal
	PriceChangePct  decimal.Decimal

	// Peak/trough tracking
	HighestPrice    decimal.Decimal
	HighestPricePct decimal.Decimal
	HighestPriceAt  time.Time

	LowestPrice    decimal.Decimal
	LowestPricePct decimal.Decimal
	LowestPriceAt  time.Time

	CreatedAt time.Time
}

// NewSignalTracking creates a new signal tracking record
func NewSignalTracking(signalID string, signal *Signal, currentPrice decimal.Decimal) *SignalTracking {
	now := time.Now()
	hoursElapsed := decimal.NewFromFloat(time.Since(signal.GeneratedAt).Hours())
	priceChangePct := signal.CalculatePriceChange(currentPrice)

	return &SignalTracking{
		SignalID:        signalID,
		TrackedAt:       now,
		HoursElapsed:    hoursElapsed,
		CurrentPrice:    currentPrice,
		PriceChangePct:  priceChangePct,
		HighestPrice:    currentPrice,
		HighestPricePct: priceChangePct,
		HighestPriceAt:  now,
		LowestPrice:     currentPrice,
		LowestPricePct:  priceChangePct,
		LowestPriceAt:   now,
		CreatedAt:       now,
	}
}

// UpdatePeakTrough updates the peak and trough prices
func (st *SignalTracking) UpdatePeakTrough(currentPrice, priceChangePct decimal.Decimal) {
	now := time.Now()

	// Update highest if current is higher
	if priceChangePct.GreaterThan(st.HighestPricePct) {
		st.HighestPrice = currentPrice
		st.HighestPricePct = priceChangePct
		st.HighestPriceAt = now
	}

	// Update lowest if current is lower
	if priceChangePct.LessThan(st.LowestPricePct) {
		st.LowestPrice = currentPrice
		st.LowestPricePct = priceChangePct
		st.LowestPriceAt = now
	}
}

// Validate validates the signal tracking
func (st *SignalTracking) Validate() error {
	if st.SignalID == "" {
		return fmt.Errorf("signal_id is required")
	}

	if st.CurrentPrice.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("current_price must be positive")
	}

	return nil
}

// SignalOutcome represents the final outcome of a signal
type SignalOutcome struct {
	ID       int64
	SignalID string

	// Outcome details
	Outcome string // PROFIT, LOSS, NEUTRAL, TIMEOUT

	// Performance metrics
	MaxFavorableMovePct decimal.Decimal
	MaxAdverseMovePct   decimal.Decimal
	FinalPriceChangePct decimal.Decimal

	// Timing
	HoursToPeak    *int
	HoursToTrough  *int
	TotalTrackingHours int

	// Additional metrics
	ProfitTargetHit bool
	StopLossHit     bool

	ClosedAt  time.Time
	CreatedAt time.Time
}

// OutcomeType represents the type of outcome
type OutcomeType string

const (
	OutcomeProfit  OutcomeType = "PROFIT"
	OutcomeLoss    OutcomeType = "LOSS"
	OutcomeNeutral OutcomeType = "NEUTRAL"
	OutcomeTimeout OutcomeType = "TIMEOUT"
)

// NewSignalOutcome creates a new signal outcome from tracking data
func NewSignalOutcome(
	signalID string,
	signal *Signal,
	finalTracking *SignalTracking,
	profitTargetPct, stopLossPct decimal.Decimal,
) *SignalOutcome {
	now := time.Now()

	// Determine outcome
	outcome := determineOutcome(finalTracking.PriceChangePct, profitTargetPct, stopLossPct)

	// Calculate hours to peak/trough
	var hoursToPeak, hoursToTrough *int
	if !finalTracking.HighestPriceAt.IsZero() {
		hours := int(finalTracking.HighestPriceAt.Sub(signal.GeneratedAt).Hours())
		hoursToPeak = &hours
	}
	if !finalTracking.LowestPriceAt.IsZero() {
		hours := int(finalTracking.LowestPriceAt.Sub(signal.GeneratedAt).Hours())
		hoursToTrough = &hours
	}

	return &SignalOutcome{
		SignalID:            signalID,
		Outcome:             string(outcome),
		MaxFavorableMovePct: finalTracking.HighestPricePct,
		MaxAdverseMovePct:   finalTracking.LowestPricePct,
		FinalPriceChangePct: finalTracking.PriceChangePct,
		HoursToPeak:         hoursToPeak,
		HoursToTrough:       hoursToTrough,
		TotalTrackingHours:  int(finalTracking.HoursElapsed.IntPart()),
		ProfitTargetHit:     finalTracking.HighestPricePct.GreaterThanOrEqual(profitTargetPct),
		StopLossHit:         finalTracking.LowestPricePct.LessThanOrEqual(stopLossPct.Neg()),
		ClosedAt:            now,
		CreatedAt:           now,
	}
}

// determineOutcome determines the outcome based on price change
func determineOutcome(priceChangePct, profitTargetPct, stopLossPct decimal.Decimal) OutcomeType {
	if priceChangePct.GreaterThanOrEqual(profitTargetPct) {
		return OutcomeProfit
	}

	if priceChangePct.LessThanOrEqual(stopLossPct.Neg()) {
		return OutcomeLoss
	}

	if priceChangePct.GreaterThan(decimal.Zero) {
		return OutcomeProfit
	}

	if priceChangePct.LessThan(decimal.Zero) {
		return OutcomeLoss
	}

	return OutcomeNeutral
}

// Validate validates the signal outcome
func (so *SignalOutcome) Validate() error {
	if so.SignalID == "" {
		return fmt.Errorf("signal_id is required")
	}

	validOutcomes := map[string]bool{
		string(OutcomeProfit):  true,
		string(OutcomeLoss):    true,
		string(OutcomeNeutral): true,
		string(OutcomeTimeout): true,
	}

	if !validOutcomes[so.Outcome] {
		return fmt.Errorf("invalid outcome: %s", so.Outcome)
	}

	return nil
}

// IsProfit returns true if the outcome is profitable
func (so *SignalOutcome) IsProfit() bool {
	return so.Outcome == string(OutcomeProfit)
}

// IsLoss returns true if the outcome is a loss
func (so *SignalOutcome) IsLoss() bool {
	return so.Outcome == string(OutcomeLoss)
}
