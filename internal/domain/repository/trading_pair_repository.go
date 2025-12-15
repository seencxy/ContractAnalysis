package repository

import (
	"context"
)

// TradingPair represents a trading pair entity
type TradingPair struct {
	ID       int64
	Symbol   string
	IsActive bool
}

// TradingPairRepository defines the interface for trading pair storage
type TradingPairRepository interface {
	// Create creates a new trading pair
	Create(ctx context.Context, pair *TradingPair) error

	// CreateBatch creates multiple trading pairs in a batch
	CreateBatch(ctx context.Context, pairs []*TradingPair) error

	// Update updates an existing trading pair
	Update(ctx context.Context, pair *TradingPair) error

	// GetBySymbol retrieves a trading pair by symbol
	GetBySymbol(ctx context.Context, symbol string) (*TradingPair, error)

	// GetAll retrieves all trading pairs
	GetAll(ctx context.Context) ([]*TradingPair, error)

	// GetActive retrieves all active trading pairs
	GetActive(ctx context.Context) ([]*TradingPair, error)

	// SetActive sets the active status of a trading pair
	SetActive(ctx context.Context, symbol string, isActive bool) error

	// Exists checks if a trading pair exists
	Exists(ctx context.Context, symbol string) (bool, error)
}
