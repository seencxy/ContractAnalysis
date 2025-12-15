package dto

// SignalListRequest represents request parameters for signal list
type SignalListRequest struct {
	FilterRequest
	TimeRangeRequest
	Symbol   string `form:"symbol"`
	Status   string `form:"status" binding:"omitempty,oneof=PENDING CONFIRMED TRACKING CLOSED INVALIDATED"`
	Type     string `form:"type" binding:"omitempty,oneof=LONG SHORT"`
	Strategy string `form:"strategy"`
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
