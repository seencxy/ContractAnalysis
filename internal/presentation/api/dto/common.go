package dto

import "time"

// TimeRangeRequest represents a time range filter
type TimeRangeRequest struct {
	StartTime *time.Time `form:"start_time"`
	EndTime   *time.Time `form:"end_time"`
}

// FilterRequest represents common filter parameters
type FilterRequest struct {
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	Sort  string `form:"sort"`
	Order string `form:"order"` // asc or desc
}

// PeriodRequest represents a time period filter
type PeriodRequest struct {
	Period string `form:"period" binding:"omitempty,oneof=24h 7d 30d all"`
}
