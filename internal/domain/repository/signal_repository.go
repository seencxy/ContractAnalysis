package repository

import (
	"context"
	"time"

	"ContractAnalysis/internal/domain/entity"
)

// SignalRepository defines the interface for signal storage
type SignalRepository interface {
	// Create creates a new signal
	Create(ctx context.Context, signal *entity.Signal) error

	// Update updates an existing signal
	Update(ctx context.Context, signal *entity.Signal) error

	// GetByID retrieves a signal by its UUID
	GetByID(ctx context.Context, signalID string) (*entity.Signal, error)

	// GetBySymbol retrieves signals for a symbol
	GetBySymbol(ctx context.Context, symbol string, limit int) ([]*entity.Signal, error)

	// GetByStatus retrieves signals with a specific status
	GetByStatus(ctx context.Context, status entity.SignalStatus, limit int) ([]*entity.Signal, error)

	// GetAll retrieves all signals (for statistics calculation)
	GetAll(ctx context.Context) ([]*entity.Signal, error)

	// GetActiveSignals retrieves all active signals (PENDING, CONFIRMED, TRACKING)
	GetActiveSignals(ctx context.Context) ([]*entity.Signal, error)

	// GetPendingSignals retrieves all pending signals
	GetPendingSignals(ctx context.Context) ([]*entity.Signal, error)

	// GetConfirmedSignals retrieves all confirmed signals
	GetConfirmedSignals(ctx context.Context) ([]*entity.Signal, error)

	// GetTrackingSignals retrieves all signals being tracked
	GetTrackingSignals(ctx context.Context) ([]*entity.Signal, error)

	// GetRecentSignalsBySymbol retrieves recent signals for a symbol within a time window
	GetRecentSignalsBySymbol(ctx context.Context, symbol string, since time.Time) ([]*entity.Signal, error)

	// CountActiveSignalsBySymbol counts active signals for a symbol
	CountActiveSignalsBySymbol(ctx context.Context, symbol string) (int, error)

	// GetSignalsInTimeRange retrieves signals generated within a time range
	GetSignalsInTimeRange(ctx context.Context, start, end time.Time) ([]*entity.Signal, error)

	// GetSignalsByStrategy retrieves signals for a specific strategy
	GetSignalsByStrategy(ctx context.Context, strategyName string, limit int) ([]*entity.Signal, error)

	// CreateTracking creates a new signal tracking record
	CreateTracking(ctx context.Context, tracking *entity.SignalTracking) error

	// GetLatestTracking retrieves the latest tracking record for a signal
	GetLatestTracking(ctx context.Context, signalID string) (*entity.SignalTracking, error)

	// GetAllTracking retrieves all tracking records for a signal
	GetAllTracking(ctx context.Context, signalID string) ([]*entity.SignalTracking, error)

	// CreateOutcome creates a new signal outcome
	CreateOutcome(ctx context.Context, outcome *entity.SignalOutcome) error

	// GetOutcome retrieves the outcome for a signal
	GetOutcome(ctx context.Context, signalID string) (*entity.SignalOutcome, error)

	// GetOutcomesBySignalIDs retrieves outcomes for multiple signals
	GetOutcomesBySignalIDs(ctx context.Context, signalIDs []string) (map[string]*entity.SignalOutcome, error)

	// GetOutcomesByTimeRange retrieves outcomes within a time range
	GetOutcomesByTimeRange(ctx context.Context, start, end time.Time) ([]*entity.SignalOutcome, error)

	// GetOutcomesByStrategy retrieves outcomes for a specific strategy
	GetOutcomesByStrategy(ctx context.Context, strategyName string, start, end time.Time) ([]*entity.SignalOutcome, error)

	// Kline tracking methods

	// CreateKlineTracking creates a new kline tracking record
	CreateKlineTracking(ctx context.Context, tracking *entity.SignalKlineTracking) error

	// GetKlineTrackingBySignal retrieves all kline tracking records for a signal
	GetKlineTrackingBySignal(ctx context.Context, signalID string) ([]*entity.SignalKlineTracking, error)

	// GetLatestKlineTracking retrieves the latest kline tracking record for a signal
	GetLatestKlineTracking(ctx context.Context, signalID string) (*entity.SignalKlineTracking, error)

	// GetKlineTrackingInTimeRange retrieves kline tracking records within a time range
	GetKlineTrackingInTimeRange(ctx context.Context, start, end time.Time) ([]*entity.SignalKlineTracking, error)
}
