package service

import (
	"context"

	"ContractAnalysis/internal/domain/entity"
)

// Strategy defines the interface for all trading strategies
type Strategy interface {
	// Name returns the strategy name
	Name() string

	// IsEnabled returns whether the strategy is enabled
	IsEnabled() bool

	// Analyze analyzes market data and generates signals
	// Takes a list of recent market data (ordered by time, newest first)
	// Returns a list of generated signals
	Analyze(ctx context.Context, recentData []*entity.MarketData) ([]*entity.Signal, error)

	// ShouldGenerateSignal checks if conditions are met to generate a signal
	ShouldGenerateSignal(ctx context.Context, data *entity.MarketData) (bool, string, error)

	// GetConfirmationHours returns the required confirmation period in hours
	GetConfirmationHours() int

	// GetTrackingHours returns the tracking period in hours
	GetTrackingHours() int

	// GetProfitTargetPct returns the profit target percentage
	GetProfitTargetPct() float64

	// GetStopLossPct returns the stop loss percentage
	GetStopLossPct() float64
}

// StrategyConfig represents common strategy configuration
type StrategyConfig struct {
	Name              string
	Enabled           bool
	ConfirmationHours int
	TrackingHours     int
	ProfitTargetPct   float64
	StopLossPct       float64
}

// BaseStrategy provides common functionality for all strategies
type BaseStrategy struct {
	config StrategyConfig
}

// NewBaseStrategy creates a new base strategy
func NewBaseStrategy(config StrategyConfig) *BaseStrategy {
	return &BaseStrategy{
		config: config,
	}
}

// Name returns the strategy name
func (s *BaseStrategy) Name() string {
	return s.config.Name
}

// IsEnabled returns whether the strategy is enabled
func (s *BaseStrategy) IsEnabled() bool {
	return s.config.Enabled
}

// GetConfirmationHours returns the confirmation period in hours
func (s *BaseStrategy) GetConfirmationHours() int {
	return s.config.ConfirmationHours
}

// GetTrackingHours returns the tracking period in hours
func (s *BaseStrategy) GetTrackingHours() int {
	return s.config.TrackingHours
}

// GetProfitTargetPct returns the profit target percentage
func (s *BaseStrategy) GetProfitTargetPct() float64 {
	return s.config.ProfitTargetPct
}

// GetStopLossPct returns the stop loss percentage
func (s *BaseStrategy) GetStopLossPct() float64 {
	return s.config.StopLossPct
}
