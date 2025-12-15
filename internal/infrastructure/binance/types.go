package binance

import (
	"time"
)

// SymbolInfo represents basic trading pair information
type SymbolInfo struct {
	Symbol     string `json:"symbol"`
	Status     string `json:"status"`
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
}

// ExchangeInfo represents the exchange information response
type ExchangeInfo struct {
	Symbols []SymbolInfo `json:"symbols"`
}

// LongShortRatio represents the long/short ratio data
type LongShortRatio struct {
	Symbol         string  `json:"symbol"`
	LongAccount    float64 `json:"longAccount,string"`
	ShortAccount   float64 `json:"shortAccount,string"`
	LongShortRatio float64 `json:"longShortRatio,string"`
	Timestamp      int64   `json:"timestamp"`
}

// TopLongShortAccountRatio represents top trader long/short account ratio
type TopLongShortAccountRatio struct {
	Symbol         string  `json:"symbol"`
	LongAccount    float64 `json:"longAccount,string"`
	ShortAccount   float64 `json:"shortAccount,string"`
	LongShortRatio float64 `json:"longShortRatio,string"`
	Timestamp      int64   `json:"timestamp"`
}

// TopLongShortPositionRatio represents top trader long/short position ratio
type TopLongShortPositionRatio struct {
	Symbol         string  `json:"symbol"`
	LongAccount    float64 `json:"longAccount,string"`    // API returns longAccount, not longPosition
	ShortAccount   float64 `json:"shortAccount,string"`   // API returns shortAccount, not shortPosition
	LongShortRatio float64 `json:"longShortRatio,string"`
	Timestamp      int64   `json:"timestamp"`
}

// GlobalLongShortAccountRatio represents global long/short account ratio
type GlobalLongShortAccountRatio struct {
	Symbol         string  `json:"symbol"`
	LongAccount    float64 `json:"longAccount,string"`
	ShortAccount   float64 `json:"shortAccount,string"`
	LongShortRatio float64 `json:"longShortRatio,string"`
	Timestamp      int64   `json:"timestamp"`
}

// TickerPrice represents the current price for a symbol
type TickerPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
	Time   int64   `json:"time"`
}

// Ticker24hr represents 24-hour ticker statistics
type Ticker24hr struct {
	Symbol             string  `json:"symbol"`
	PriceChange        float64 `json:"priceChange,string"`
	PriceChangePercent float64 `json:"priceChangePercent,string"`
	LastPrice          float64 `json:"lastPrice,string"`
	Volume             float64 `json:"volume,string"`
	QuoteVolume        float64 `json:"quoteVolume,string"`
	OpenTime           int64   `json:"openTime"`
	CloseTime          int64   `json:"closeTime"`
	Count              int64   `json:"count"`
}

// MarketData represents collected market data for a symbol
type MarketData struct {
	Symbol    string
	Timestamp time.Time

	// Long/Short ratios (as percentages 0-100)
	LongAccountRatio  float64
	ShortAccountRatio float64

	LongPositionRatio  float64
	ShortPositionRatio float64

	// Trader counts (if available)
	LongTraderCount  int
	ShortTraderCount int

	// Price and volume
	Price     float64
	Volume24h float64
}

// RateLimitInfo represents rate limit information
type RateLimitInfo struct {
	RequestsRemaining int
	WeightRemaining   int
	ResetTime         time.Time
}

// Kline represents a kline/candlestick data point (alias for entity.Kline to avoid import cycle)
type Kline struct {
	OpenTime    time.Time
	CloseTime   time.Time
	Open        float64
	High        float64
	Low         float64
	Close       float64
	Volume      float64
	QuoteVolume float64
}
