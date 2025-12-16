package repository

import (
	"context"
	"time"

	"ContractAnalysis/internal/domain/entity"
)

// KlineRepository defines the interface for accessing kline data
type KlineRepository interface {
	// GetKlines retrieves kline data for a symbol
	GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]*entity.Kline, error)

	// GetKlinesSince retrieves kline data since a specific time
	GetKlinesSince(ctx context.Context, symbol string, interval string, startTime time.Time) ([]*entity.Kline, error)
}
