package service

import (
	"ContractAnalysis/internal/domain/entity"

	"github.com/shopspring/decimal"
)

// PatternAnalyzer provides methods for candlestick pattern detection
type PatternAnalyzer struct{}

// NewPatternAnalyzer creates a new pattern analyzer
func NewPatternAnalyzer() *PatternAnalyzer {
	return &PatternAnalyzer{}
}

// IsShootingStar checks if a candle matches the Shooting Star pattern
// A Shooting Star is a bearish reversal pattern formed by a single candle
// Characteristics:
// 1. Small body at the lower end
// 2. Long upper wick (at least 2x the body length)
// 3. Very short or no lower wick
// 4. Occurs after an uptrend (context check needed outside this function)
func (p *PatternAnalyzer) IsShootingStar(k *entity.Kline) bool {
	open := k.Open
	close := k.Close
	high := k.High
	low := k.Low

	// Calculate body and wicks
	bodySize := open.Sub(close).Abs()
	upperWick := high.Sub(decimal.Max(open, close))
	lowerWick := decimal.Min(open, close).Sub(low)

	// Check 1: Small body (relative to total range)
	totalRange := high.Sub(low)
	if totalRange.IsZero() {
		return false
	}
	// Body should be in the lower third? Or just small.
	// Let's say body is less than 30% of total range (configurable, but standard is small)
	if bodySize.Div(totalRange).GreaterThan(decimal.NewFromFloat(0.4)) {
		return false
	}

	// Check 2: Long upper wick (>= 2 * body size)
	// Handle zero body size (doji)
	minBody := bodySize
	if minBody.IsZero() {
		// Use a tiny epsilon or just compare to range
		minBody = totalRange.Mul(decimal.NewFromFloat(0.01))
	}

	if upperWick.LessThan(minBody.Mul(decimal.NewFromInt(2))) {
		return false
	}

	// Check 3: Short lower wick (<= body size or very small)
	// Strictly, it should be no lower wick, but real markets have noise.
	// Let's say lower wick <= body size
	if lowerWick.GreaterThan(minBody) {
		return false
	}

	return true
}

// IsBearishEngulfing checks for a Bearish Engulfing pattern
// Characteristics:
// 1. Previous candle is bullish (Green)
// 2. Current candle is bearish (Red)
// 3. Current body completely engulfs previous body (Open >= Prev Close and Close <= Prev Open)
func (p *PatternAnalyzer) IsBearishEngulfing(current, previous *entity.Kline) bool {
	// 1. Previous candle bullish
	if previous.Close.LessThanOrEqual(previous.Open) {
		return false
	}

	// 2. Current candle bearish
	if current.Close.GreaterThanOrEqual(current.Open) {
		return false
	}

	// 3. Engulfing check
	// In crypto, Open of current usually equals Close of previous, so we use >= and <=
	// Current Open should be >= Previous Close
	// Current Close should be <= Previous Open
	if current.Open.LessThan(previous.Close) {
		return false
	}
	if current.Close.GreaterThan(previous.Open) {
		return false
	}

	return true
}

// IsSwingFailurePattern checks for a Swing Failure Pattern (SFP)
// Returns true if the triggerCandle swept the resistanceHigh but closed below it
func (p *PatternAnalyzer) IsSwingFailurePattern(triggerCandle *entity.Kline, resistanceHigh decimal.Decimal) bool {
	// 1. High must be above resistance (Sweep)
	// 2. Close must be below resistance (Failure)
	return triggerCandle.High.GreaterThan(resistanceHigh) && triggerCandle.Close.LessThan(resistanceHigh)
}
