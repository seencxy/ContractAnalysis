package entity

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// SignalKlineTracking represents kline-based tracking data for a signal
type SignalKlineTracking struct {
	ID       int64
	SignalID string

	// K-line time information
	KlineOpenTime    time.Time
	KlineCloseTime   time.Time
	HoursSinceSignal decimal.Decimal

	// OHLCV data
	OpenPrice   decimal.Decimal
	HighPrice   decimal.Decimal
	LowPrice    decimal.Decimal
	ClosePrice  decimal.Decimal
	Volume      decimal.Decimal
	QuoteVolume decimal.Decimal

	// Price change percentages (relative to signal price)
	OpenChangePct  decimal.Decimal
	HighChangePct  decimal.Decimal
	LowChangePct   decimal.Decimal
	CloseChangePct decimal.Decimal

	// Hourly return
	HourlyReturnPct decimal.Decimal

	// Theoretical maximum profit/loss
	MaxPotentialProfitPct decimal.Decimal
	MaxPotentialLossPct   decimal.Decimal

	// Profitability flags
	IsProfitableAtHigh  bool
	IsProfitableAtClose bool

	CreatedAt time.Time
}

// Kline represents a kline (candlestick) data point
type Kline struct {
	OpenTime    time.Time
	CloseTime   time.Time
	Open        decimal.Decimal
	High        decimal.Decimal
	Low         decimal.Decimal
	Close       decimal.Decimal
	Volume      decimal.Decimal
	QuoteVolume decimal.Decimal
}

// NewSignalKlineTracking creates a new signal kline tracking record
func NewSignalKlineTracking(signalID string, signal *Signal, kline *Kline) *SignalKlineTracking {
	now := time.Now()

	// Calculate price changes relative to signal price (considering LONG/SHORT direction)
	openChangePct := calculatePriceChange(signal, kline.Open)
	highChangePct := calculatePriceChange(signal, kline.High)
	lowChangePct := calculatePriceChange(signal, kline.Low)
	closeChangePct := calculatePriceChange(signal, kline.Close)

	// Calculate hourly return: (close - open) / open * 100
	hourlyReturn := decimal.Zero
	if !kline.Open.IsZero() {
		hourlyReturn = kline.Close.Sub(kline.Open).
			Div(kline.Open).
			Mul(decimal.NewFromInt(100))
	}

	// Calculate hours since signal generated
	hoursSince := decimal.NewFromFloat(kline.OpenTime.Sub(signal.GeneratedAt).Hours())

	return &SignalKlineTracking{
		SignalID:              signalID,
		KlineOpenTime:         kline.OpenTime,
		KlineCloseTime:        kline.CloseTime,
		HoursSinceSignal:      hoursSince,
		OpenPrice:             kline.Open,
		HighPrice:             kline.High,
		LowPrice:              kline.Low,
		ClosePrice:            kline.Close,
		Volume:                kline.Volume,
		QuoteVolume:           kline.QuoteVolume,
		OpenChangePct:         openChangePct,
		HighChangePct:         highChangePct,
		LowChangePct:          lowChangePct,
		CloseChangePct:        closeChangePct,
		HourlyReturnPct:       hourlyReturn,
		MaxPotentialProfitPct: highChangePct,
		MaxPotentialLossPct:   lowChangePct,
		IsProfitableAtHigh:    highChangePct.GreaterThan(decimal.Zero),
		IsProfitableAtClose:   closeChangePct.GreaterThan(decimal.Zero),
		CreatedAt:             now,
	}
}

// calculatePriceChange calculates price change percentage considering signal direction
func calculatePriceChange(signal *Signal, currentPrice decimal.Decimal) decimal.Decimal {
	if signal.PriceAtSignal.IsZero() {
		return decimal.Zero
	}

	// Calculate percentage change: (current - signal) / signal * 100
	change := currentPrice.Sub(signal.PriceAtSignal).
		Div(signal.PriceAtSignal).
		Mul(decimal.NewFromInt(100))

	// For SHORT signals, invert the change (price decrease becomes profit)
	if signal.Type == SignalTypeShort {
		return change.Neg()
	}

	return change
}

// Validate validates the signal kline tracking record
func (kt *SignalKlineTracking) Validate() error {
	if kt.SignalID == "" {
		return fmt.Errorf("signal_id is required")
	}

	if kt.HighPrice.LessThan(kt.LowPrice) {
		return fmt.Errorf("high price cannot be less than low price")
	}

	if kt.HighPrice.LessThan(kt.OpenPrice) {
		return fmt.Errorf("high price cannot be less than open price")
	}

	if kt.HighPrice.LessThan(kt.ClosePrice) {
		return fmt.Errorf("high price cannot be less than close price")
	}

	if kt.LowPrice.GreaterThan(kt.OpenPrice) {
		return fmt.Errorf("low price cannot be greater than open price")
	}

	if kt.LowPrice.GreaterThan(kt.ClosePrice) {
		return fmt.Errorf("low price cannot be greater than close price")
	}

	return nil
}
