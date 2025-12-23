package entity

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// MarketData represents market data for a trading pair at a specific time
type MarketData struct {
	ID        int64
	Symbol    string
	Timestamp time.Time

	// Long/Short ratio by account count
	LongAccountRatio  decimal.Decimal
	ShortAccountRatio decimal.Decimal

	// Long/Short ratio by position size (whale positions)
	LongPositionRatio  decimal.Decimal
	ShortPositionRatio decimal.Decimal

	// Data quality indicators
	PositionRatioAvailable bool // Whether position ratio data is available from API
	DataQualityScore       int  // Data quality score 0-100

	// Price and volume
	Price        decimal.Decimal
	Volume24h    decimal.Decimal // Optional
	OpenInterest decimal.Decimal // Open Interest in USDT
	FundingRate  decimal.Decimal // Current Funding Rate

	CreatedAt time.Time
}

// Validate validates the market data
func (m *MarketData) Validate() error {
	if m.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}

	if m.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	// Tolerance for floating point precision errors
	tolerance := decimal.NewFromFloat(0.01)
	hundred := decimal.NewFromInt(100)

	// Validate account ratios sum to 100 (with tolerance for floating point errors)
	accountSum := m.LongAccountRatio.Add(m.ShortAccountRatio)
	accountDiff := accountSum.Sub(hundred).Abs()
	if accountDiff.GreaterThan(tolerance) {
		return fmt.Errorf("account ratios must sum to 100 (±0.01), got: %s", accountSum.String())
	}

	// Validate position ratios sum to 100 (with tolerance)
	// Allow position ratios to be 0 (some pairs don't have this data from Binance)
	positionSum := m.LongPositionRatio.Add(m.ShortPositionRatio)
	if !positionSum.IsZero() {
		positionDiff := positionSum.Sub(hundred).Abs()
		if positionDiff.GreaterThan(tolerance) {
			return fmt.Errorf("position ratios must sum to 100 (±0.01), got: %s", positionSum.String())
		}
	}

	// Validate ratios are between 0 and 100
	if m.LongAccountRatio.LessThan(decimal.Zero) || m.LongAccountRatio.GreaterThan(hundred) {
		return fmt.Errorf("long account ratio must be between 0 and 100")
	}

	if m.ShortAccountRatio.LessThan(decimal.Zero) || m.ShortAccountRatio.GreaterThan(hundred) {
		return fmt.Errorf("short account ratio must be between 0 and 100")
	}

	if m.Price.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("price must be positive")
	}

	if m.OpenInterest.LessThan(decimal.Zero) {
		return fmt.Errorf("open interest must be non-negative")
	}

	// Validate timestamp freshness
	now := time.Now()

	// Don't allow future timestamps (with 5-minute clock skew tolerance)
	if m.Timestamp.After(now.Add(5 * time.Minute)) {
		return fmt.Errorf("future timestamp detected: %v (now: %v)", m.Timestamp, now)
	}

	// Don't allow stale data (older than 1 hour)
	if m.Timestamp.Before(now.Add(-1 * time.Hour)) {
		return fmt.Errorf("stale data detected: %v (now: %v)", m.Timestamp, now)
	}

	// Validate data quality score
	if m.DataQualityScore < 0 || m.DataQualityScore > 100 {
		return fmt.Errorf("data quality score must be between 0 and 100, got: %d", m.DataQualityScore)
	}

	return nil
}

// CalculateRatioDifference calculates the difference between long and short ratios
func (m *MarketData) CalculateRatioDifference() decimal.Decimal {
	return m.LongAccountRatio.Sub(m.ShortAccountRatio).Abs()
}

// GetDominantDirection returns the dominant direction based on account ratio
func (m *MarketData) GetDominantDirection() string {
	if m.LongAccountRatio.GreaterThan(m.ShortAccountRatio) {
		return "LONG"
	}
	return "SHORT"
}

// GetDominantRatio returns the larger of the two ratios
func (m *MarketData) GetDominantRatio() decimal.Decimal {
	if m.LongAccountRatio.GreaterThan(m.ShortAccountRatio) {
		return m.LongAccountRatio
	}
	return m.ShortAccountRatio
}

// GetMinorityDirection returns the minority direction (opposite of dominant)
func (m *MarketData) GetMinorityDirection() string {
	if m.GetDominantDirection() == "LONG" {
		return "SHORT"
	}
	return "LONG"
}

// HasDivergence checks if there's divergence between account ratio and position ratio
// Returns true if the dominant direction differs between account count and position size
func (m *MarketData) HasDivergence() bool {
	accountDominant := m.LongAccountRatio.GreaterThan(m.ShortAccountRatio)
	positionDominant := m.LongPositionRatio.GreaterThan(m.ShortPositionRatio)
	return accountDominant != positionDominant
}

// CalculateDivergence calculates the divergence percentage between account and position ratios
func (m *MarketData) CalculateDivergence() decimal.Decimal {
	accountDiff := m.LongAccountRatio.Sub(m.ShortAccountRatio)
	positionDiff := m.LongPositionRatio.Sub(m.ShortPositionRatio)
	return accountDiff.Sub(positionDiff).Abs()
}

// IsAccountRatioExtreme checks if the account ratio is extreme (one side dominates)
func (m *MarketData) IsAccountRatioExtreme(threshold decimal.Decimal) bool {
	return m.GetDominantRatio().GreaterThanOrEqual(threshold)
}

// IsPositionRatioExtreme checks if the position ratio is extreme
func (m *MarketData) IsPositionRatioExtreme(threshold decimal.Decimal) bool {
	dominantPosition := m.LongPositionRatio
	if m.ShortPositionRatio.GreaterThan(m.LongPositionRatio) {
		dominantPosition = m.ShortPositionRatio
	}
	return dominantPosition.GreaterThanOrEqual(threshold)
}

// GetWhaleDirection returns the direction whales are taking based on position ratio
func (m *MarketData) GetWhaleDirection() string {
	if m.LongPositionRatio.GreaterThan(m.ShortPositionRatio) {
		return "LONG"
	}
	return "SHORT"
}

// IsValid is a convenience method that calls Validate and returns a bool
func (m *MarketData) IsValid() bool {
	return m.Validate() == nil
}

// String returns a string representation of the market data
func (m *MarketData) String() string {
	return fmt.Sprintf("MarketData{Symbol: %s, Time: %s, LongAcct: %s%%, ShortAcct: %s%%, LongPos: %s%%, ShortPos: %s%%, Price: %s, OI: %s, FR: %s}",
		m.Symbol,
		m.Timestamp.Format(time.RFC3339),
		m.LongAccountRatio.String(),
		m.ShortAccountRatio.String(),
		m.LongPositionRatio.String(),
		m.ShortPositionRatio.String(),
		m.Price.String(),
		m.OpenInterest.String(),
		m.FundingRate.String(),
	)
}
