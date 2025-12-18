package dto

// SignalListRequest represents request parameters for signal list
type SignalListRequest struct {
	FilterRequest
	TimeRangeRequest
	Symbol       string `form:"symbol"`
	Status       string `form:"status" binding:"omitempty,oneof=PENDING CONFIRMED TRACKING CLOSED INVALIDATED"`
	Type         string `form:"type" binding:"omitempty,oneof=LONG SHORT"`
	StrategyName string `form:"strategy_name"`
}

// StatisticsRequest represents request parameters for statistics
type StatisticsRequest struct {
	PeriodRequest
	StrategyName string `form:"strategy"`
	Symbol       string `form:"symbol"`
}

// PairListRequest represents request parameters for trading pair list
type PairListRequest struct {
	FilterRequest
	IsActive *bool `form:"is_active"`
}

// MarketDataRequest represents request parameters for market data
type MarketDataRequest struct {
	TimeRangeRequest
	Symbol string `form:"symbol" binding:"required"`
	Limit  int    `form:"limit"`
}

// StatisticsHistoryRequest represents request for historical statistics
type StatisticsHistoryRequest struct {
	TimeRangeRequest
	StrategyName string `form:"strategy"`
	Symbol       string `form:"symbol"`
}

// StrategyCompareRequest represents request parameters for strategy comparison
type StrategyCompareRequest struct {
	StrategyNames []string `form:"strategies" binding:"required,min=2,max=5"` // 2-5 strategies
	Period        string   `form:"period" binding:"required,oneof=24h 7d 30d all"`
	Symbols       []string `form:"symbols"` // Optional: filter by specific symbols
}
