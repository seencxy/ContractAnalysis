package repository

import (
	"context"
	"time"

	"ContractAnalysis/internal/domain/entity"
)

// MarketDataRepository defines the interface for market data storage
type MarketDataRepository interface {
	// Create creates a new market data record
	Create(ctx context.Context, data *entity.MarketData) error

	// CreateBatch creates multiple market data records in a batch
	CreateBatch(ctx context.Context, dataList []*entity.MarketData) error

	// GetBySymbol retrieves market data for a symbol within a time range
	GetBySymbol(ctx context.Context, symbol string, start, end time.Time) ([]*entity.MarketData, error)

	// GetLatestBySymbol retrieves the latest market data for a symbol
	GetLatestBySymbol(ctx context.Context, symbol string) (*entity.MarketData, error)

	// GetLatestForAllSymbols retrieves the latest market data for all symbols
	GetLatestForAllSymbols(ctx context.Context) ([]*entity.MarketData, error)

	// GetRecentBySymbol retrieves the most recent N records for a symbol
	GetRecentBySymbol(ctx context.Context, symbol string, limit int) ([]*entity.MarketData, error)

	// Delete deletes market data older than the specified time
	DeleteOlderThan(ctx context.Context, before time.Time) error

	// Count returns the total number of market data records
	Count(ctx context.Context) (int64, error)
}
