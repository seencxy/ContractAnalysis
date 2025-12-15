package dto

import "time"

// SignalResponse represents a signal in API response
type SignalResponse struct {
	SignalID           string  `json:"signal_id"`
	Symbol             string  `json:"symbol"`
	Type               string  `json:"type"`
	StrategyName       string  `json:"strategy_name"`
	GeneratedAt        string  `json:"generated_at"`
	PriceAtSignal      string  `json:"price_at_signal"`
	LongAccountRatio   string  `json:"long_account_ratio"`
	ShortAccountRatio  string  `json:"short_account_ratio"`
	LongPositionRatio  string  `json:"long_position_ratio"`
	ShortPositionRatio string  `json:"short_position_ratio"`
	LongTraderCount    int     `json:"long_trader_count"`
	ShortTraderCount   int     `json:"short_trader_count"`
	Status             string  `json:"status"`
	IsConfirmed        bool    `json:"is_confirmed"`
	ConfirmedAt        *string `json:"confirmed_at,omitempty"`
	Reason             string  `json:"reason,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`

	// Final outcome (only for CLOSED signals)
	FinalPnlPct        *string `json:"final_pnl_pct,omitempty"`        // 最终盈亏百分比
	Outcome            *string `json:"outcome,omitempty"`              // PROFIT, LOSS, NEUTRAL
	TotalTrackingHours *int    `json:"total_tracking_hours,omitempty"` // 总追踪小时数
	ClosedAt           *string `json:"closed_at,omitempty"`            // 关闭时间（仅已关闭信号）
}

// SignalTrackingResponse represents signal tracking data
type SignalTrackingResponse struct {
	ID                int64   `json:"id"`
	SignalID          string  `json:"signal_id"`
	TrackedAt         string  `json:"tracked_at"`
	CurrentPrice      string  `json:"current_price"`
	PriceChangePct    string  `json:"price_change_pct"`
	HighestPrice      *string `json:"highest_price,omitempty"`
	LowestPrice       *string `json:"lowest_price,omitempty"`
	HighestChangePct  *string `json:"highest_change_pct,omitempty"`
	LowestChangePct   *string `json:"lowest_change_pct,omitempty"`
	HoursTracked      int     `json:"hours_tracked"`
	IsProfitTargetHit bool    `json:"is_profit_target_hit"`
	IsStopLossHit     bool    `json:"is_stop_loss_hit"`
}

// SignalKlineTrackingResponse represents K-line tracking data
type SignalKlineTrackingResponse struct {
	ID                  int64   `json:"id"`
	SignalID            string  `json:"signal_id"`
	KlineOpenTime       string  `json:"kline_open_time"`
	KlineCloseTime      string  `json:"kline_close_time"`
	OpenPrice           string  `json:"open_price"`
	HighPrice           string  `json:"high_price"`
	LowPrice            string  `json:"low_price"`
	ClosePrice          string  `json:"close_price"`
	Volume              string  `json:"volume"`
	OpenChangePct       string  `json:"open_change_pct"`
	HighChangePct       string  `json:"high_change_pct"`
	LowChangePct        string  `json:"low_change_pct"`
	CloseChangePct      string  `json:"close_change_pct"`
	HourlyReturnPct     *string `json:"hourly_return_pct,omitempty"`
	IsProfitableAtHigh  bool    `json:"is_profitable_at_high"`
	IsProfitableAtClose bool    `json:"is_profitable_at_close"`
}

// StatisticsResponse represents strategy statistics
type StatisticsResponse struct {
	StrategyName string  `json:"strategy_name"`
	Symbol       *string `json:"symbol,omitempty"`
	PeriodLabel  string  `json:"period_label"`
	PeriodStart  string  `json:"period_start"`
	PeriodEnd    string  `json:"period_end"`

	// Signal counts
	TotalSignals       int `json:"total_signals"`
	ConfirmedSignals   int `json:"confirmed_signals"`
	InvalidatedSignals int `json:"invalidated_signals"`

	// Outcome counts
	ProfitableSignals int `json:"profitable_signals"`
	LosingSignals     int `json:"losing_signals"`
	NeutralSignals    int `json:"neutral_signals"`

	// Performance metrics
	WinRate         *string `json:"win_rate,omitempty"`
	AvgProfitPct    *string `json:"avg_profit_pct,omitempty"`
	AvgLossPct      *string `json:"avg_loss_pct,omitempty"`
	AvgHoldingHours *string `json:"avg_holding_hours,omitempty"`

	// Best/Worst
	BestSignalPct  *string `json:"best_signal_pct,omitempty"`
	WorstSignalPct *string `json:"worst_signal_pct,omitempty"`

	// Profit factor
	ProfitFactor *string `json:"profit_factor,omitempty"`

	// K-line metrics
	KlineTheoreticalWinRate   *string `json:"kline_theoretical_win_rate,omitempty"`
	KlineCloseWinRate         *string `json:"kline_close_win_rate,omitempty"`
	TotalKlineHours           int     `json:"total_kline_hours"`
	ProfitableKlineHoursHigh  int     `json:"profitable_kline_hours_high"`
	ProfitableKlineHoursClose int     `json:"profitable_kline_hours_close"`

	// Hourly return statistics
	AvgHourlyReturnPct *string `json:"avg_hourly_return_pct,omitempty"`
	MaxHourlyReturnPct *string `json:"max_hourly_return_pct,omitempty"`
	MinHourlyReturnPct *string `json:"min_hourly_return_pct,omitempty"`

	// Theoretical max profit/loss
	AvgMaxPotentialProfitPct *string `json:"avg_max_potential_profit_pct,omitempty"`
	AvgMaxPotentialLossPct   *string `json:"avg_max_potential_loss_pct,omitempty"`

	CalculatedAt string `json:"calculated_at"`
}

// TradingPairResponse represents a trading pair
type TradingPairResponse struct {
	Symbol    string `json:"symbol"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// MarketDataResponse represents market data
type MarketDataResponse struct {
	Symbol             string `json:"symbol"`
	Timestamp          string `json:"timestamp"`
	LongAccountRatio   string `json:"long_account_ratio"`
	ShortAccountRatio  string `json:"short_account_ratio"`
	LongPositionRatio  string `json:"long_position_ratio"`
	ShortPositionRatio string `json:"short_position_ratio"`
	LongTraderCount    int    `json:"long_trader_count"`
	ShortTraderCount   int    `json:"short_trader_count"`
	Volume24h          string `json:"volume_24h"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// SignalStatusDistribution represents signal count by status
type SignalStatusDistribution struct {
	Pending     int `json:"pending"`
	Confirmed   int `json:"confirmed"`
	Tracking    int `json:"tracking"`
	Closed      int `json:"closed"`
	Invalidated int `json:"invalidated"`
}

// OverviewStatisticsResponse represents overall statistics
type OverviewStatisticsResponse struct {
	TotalSignalsToday   int                       `json:"total_signals_today"`
	ActiveSignals       int                       `json:"active_signals"`
	OverallWinRate24h   *string                   `json:"overall_win_rate_24h,omitempty"`
	AvgReturnPct24h     *string                   `json:"avg_return_pct_24h,omitempty"`
	TopPerformingPair   string                    `json:"top_performing_pair,omitempty"`
	WorstPerformingPair string                    `json:"worst_performing_pair,omitempty"`
	StatusDistribution  *SignalStatusDistribution `json:"status_distribution,omitempty"`
}
