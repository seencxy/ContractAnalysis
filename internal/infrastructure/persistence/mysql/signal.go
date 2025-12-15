package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ContractAnalysis/internal/domain/entity"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// SignalModel represents the signals table
type SignalModel struct {
	ID                 int64           `gorm:"column:id;primaryKey;autoIncrement"`
	SignalID           string          `gorm:"column:signal_id;uniqueIndex;size:36;not null"`
	Symbol             string          `gorm:"column:symbol;size:50;not null;index:idx_symbol_status"`
	Type               string          `gorm:"column:signal_type;size:20;not null"`
	StrategyName       string          `gorm:"column:strategy_name;size:50;not null;index"`
	GeneratedAt        time.Time       `gorm:"column:generated_at;not null;index:idx_status_generated"`
	PriceAtSignal      decimal.Decimal `gorm:"column:price_at_signal;type:decimal(20,8);not null"`
	LongAccountRatio   decimal.Decimal `gorm:"column:long_account_ratio;type:decimal(10,4);not null"`
	ShortAccountRatio  decimal.Decimal `gorm:"column:short_account_ratio;type:decimal(10,4);not null"`
	LongPositionRatio  decimal.Decimal `gorm:"column:long_position_ratio;type:decimal(10,4);not null"`
	ShortPositionRatio decimal.Decimal `gorm:"column:short_position_ratio;type:decimal(10,4);not null"`
	LongTraderCount    int             `gorm:"column:long_trader_count;not null"`
	ShortTraderCount   int             `gorm:"column:short_trader_count;not null"`
	ConfirmationStart  time.Time       `gorm:"column:confirmation_start;not null"`
	ConfirmationEnd    time.Time       `gorm:"column:confirmation_end;not null"`
	IsConfirmed        bool            `gorm:"column:is_confirmed;default:false"`
	ConfirmedAt        *time.Time      `gorm:"column:confirmed_at"`
	Status             string          `gorm:"column:status;size:20;not null;index:idx_symbol_status;index:idx_status_generated"`
	Reason             string          `gorm:"column:reason;type:text"`
	ConfigSnapshot     string          `gorm:"column:config_snapshot;type:json"`
	CreatedAt          time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time       `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name
func (SignalModel) TableName() string {
	return "signals"
}

// ToEntity converts model to domain entity
func (m *SignalModel) ToEntity() (*entity.Signal, error) {
	var configSnapshot map[string]interface{}
	if m.ConfigSnapshot != "" {
		if err := json.Unmarshal([]byte(m.ConfigSnapshot), &configSnapshot); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config snapshot: %w", err)
		}
	}

	return &entity.Signal{
		ID:                 m.ID,
		SignalID:           m.SignalID,
		Symbol:             m.Symbol,
		Type:               entity.SignalType(m.Type),
		StrategyName:       m.StrategyName,
		GeneratedAt:        m.GeneratedAt,
		PriceAtSignal:      m.PriceAtSignal,
		LongAccountRatio:   m.LongAccountRatio,
		ShortAccountRatio:  m.ShortAccountRatio,
		LongPositionRatio:  m.LongPositionRatio,
		ShortPositionRatio: m.ShortPositionRatio,
		LongTraderCount:    m.LongTraderCount,
		ShortTraderCount:   m.ShortTraderCount,
		ConfirmationStart:  m.ConfirmationStart,
		ConfirmationEnd:    m.ConfirmationEnd,
		IsConfirmed:        m.IsConfirmed,
		ConfirmedAt:        m.ConfirmedAt,
		Status:             entity.SignalStatus(m.Status),
		Reason:             m.Reason,
		ConfigSnapshot:     configSnapshot,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}, nil
}

// FromEntity converts domain entity to model
func (m *SignalModel) FromEntity(entity *entity.Signal) error {
	var configSnapshotJSON string
	if entity.ConfigSnapshot != nil {
		data, err := json.Marshal(entity.ConfigSnapshot)
		if err != nil {
			return fmt.Errorf("failed to marshal config snapshot: %w", err)
		}
		configSnapshotJSON = string(data)
	}

	m.ID = entity.ID
	m.SignalID = entity.SignalID
	m.Symbol = entity.Symbol
	m.Type = string(entity.Type)
	m.StrategyName = entity.StrategyName
	m.GeneratedAt = entity.GeneratedAt
	m.PriceAtSignal = entity.PriceAtSignal
	m.LongAccountRatio = entity.LongAccountRatio
	m.ShortAccountRatio = entity.ShortAccountRatio
	m.LongPositionRatio = entity.LongPositionRatio
	m.ShortPositionRatio = entity.ShortPositionRatio
	m.LongTraderCount = entity.LongTraderCount
	m.ShortTraderCount = entity.ShortTraderCount
	m.ConfirmationStart = entity.ConfirmationStart
	m.ConfirmationEnd = entity.ConfirmationEnd
	m.IsConfirmed = entity.IsConfirmed
	m.ConfirmedAt = entity.ConfirmedAt
	m.Status = string(entity.Status)
	m.Reason = entity.Reason
	m.ConfigSnapshot = configSnapshotJSON

	return nil
}

// SignalTrackingModel represents the signal_tracking table
type SignalTrackingModel struct {
	ID              int64           `gorm:"column:id;primaryKey;autoIncrement"`
	SignalID        string          `gorm:"column:signal_id;size:36;not null;index:idx_signal_tracked"`
	TrackedAt       time.Time       `gorm:"column:tracked_at;not null;index:idx_signal_tracked"`
	HoursElapsed    decimal.Decimal `gorm:"column:hours_elapsed;type:decimal(10,2);not null"`
	CurrentPrice    decimal.Decimal `gorm:"column:current_price;type:decimal(20,8);not null"`
	PriceChangePct  decimal.Decimal `gorm:"column:price_change_pct;type:decimal(10,4);not null"`
	HighestPrice    decimal.Decimal `gorm:"column:highest_price;type:decimal(20,8);not null"`
	HighestPricePct decimal.Decimal `gorm:"column:highest_price_pct;type:decimal(10,4);not null"`
	HighestPriceAt  time.Time       `gorm:"column:highest_price_at;not null"`
	LowestPrice     decimal.Decimal `gorm:"column:lowest_price;type:decimal(20,8);not null"`
	LowestPricePct  decimal.Decimal `gorm:"column:lowest_price_pct;type:decimal(10,4);not null"`
	LowestPriceAt   time.Time       `gorm:"column:lowest_price_at;not null"`
	CreatedAt       time.Time       `gorm:"column:created_at;autoCreateTime"`
}

// TableName specifies the table name
func (SignalTrackingModel) TableName() string {
	return "signal_tracking"
}

// ToEntity converts model to domain entity
func (m *SignalTrackingModel) ToEntity() *entity.SignalTracking {
	return &entity.SignalTracking{
		ID:              m.ID,
		SignalID:        m.SignalID,
		TrackedAt:       m.TrackedAt,
		HoursElapsed:    m.HoursElapsed,
		CurrentPrice:    m.CurrentPrice,
		PriceChangePct:  m.PriceChangePct,
		HighestPrice:    m.HighestPrice,
		HighestPricePct: m.HighestPricePct,
		HighestPriceAt:  m.HighestPriceAt,
		LowestPrice:     m.LowestPrice,
		LowestPricePct:  m.LowestPricePct,
		LowestPriceAt:   m.LowestPriceAt,
		CreatedAt:       m.CreatedAt,
	}
}

// FromEntity converts domain entity to model
func (m *SignalTrackingModel) FromEntity(entity *entity.SignalTracking) {
	m.ID = entity.ID
	m.SignalID = entity.SignalID
	m.TrackedAt = entity.TrackedAt
	m.HoursElapsed = entity.HoursElapsed
	m.CurrentPrice = entity.CurrentPrice
	m.PriceChangePct = entity.PriceChangePct
	m.HighestPrice = entity.HighestPrice
	m.HighestPricePct = entity.HighestPricePct
	m.HighestPriceAt = entity.HighestPriceAt
	m.LowestPrice = entity.LowestPrice
	m.LowestPricePct = entity.LowestPricePct
	m.LowestPriceAt = entity.LowestPriceAt
}

// SignalOutcomeModel represents the signal_outcomes table
type SignalOutcomeModel struct {
	ID                  int64           `gorm:"column:id;primaryKey;autoIncrement"`
	SignalID            string          `gorm:"column:signal_id;uniqueIndex;size:36;not null"`
	Outcome             string          `gorm:"column:outcome;size:20;not null;index"`
	MaxFavorableMovePct decimal.Decimal `gorm:"column:max_favorable_move_pct;type:decimal(10,4);not null"`
	MaxAdverseMovePct   decimal.Decimal `gorm:"column:max_adverse_move_pct;type:decimal(10,4);not null"`
	FinalPriceChangePct decimal.Decimal `gorm:"column:final_price_change_pct;type:decimal(10,4);not null"`
	HoursToPeak         *int            `gorm:"column:hours_to_peak"`
	HoursToTrough       *int            `gorm:"column:hours_to_trough"`
	TotalTrackingHours  int             `gorm:"column:total_tracking_hours;not null"`
	ProfitTargetHit     bool            `gorm:"column:profit_target_hit;default:false"`
	StopLossHit         bool            `gorm:"column:stop_loss_hit;default:false"`
	ClosedAt            time.Time       `gorm:"column:closed_at;not null;index"`
	CreatedAt           time.Time       `gorm:"column:created_at;autoCreateTime"`
}

// TableName specifies the table name
func (SignalOutcomeModel) TableName() string {
	return "signal_outcomes"
}

// ToEntity converts model to domain entity
func (m *SignalOutcomeModel) ToEntity() *entity.SignalOutcome {
	return &entity.SignalOutcome{
		ID:                  m.ID,
		SignalID:            m.SignalID,
		Outcome:             m.Outcome,
		MaxFavorableMovePct: m.MaxFavorableMovePct,
		MaxAdverseMovePct:   m.MaxAdverseMovePct,
		FinalPriceChangePct: m.FinalPriceChangePct,
		HoursToPeak:         m.HoursToPeak,
		HoursToTrough:       m.HoursToTrough,
		TotalTrackingHours:  m.TotalTrackingHours,
		ProfitTargetHit:     m.ProfitTargetHit,
		StopLossHit:         m.StopLossHit,
		ClosedAt:            m.ClosedAt,
		CreatedAt:           m.CreatedAt,
	}
}

// FromEntity converts domain entity to model
func (m *SignalOutcomeModel) FromEntity(entity *entity.SignalOutcome) {
	m.ID = entity.ID
	m.SignalID = entity.SignalID
	m.Outcome = entity.Outcome
	m.MaxFavorableMovePct = entity.MaxFavorableMovePct
	m.MaxAdverseMovePct = entity.MaxAdverseMovePct
	m.FinalPriceChangePct = entity.FinalPriceChangePct
	m.HoursToPeak = entity.HoursToPeak
	m.HoursToTrough = entity.HoursToTrough
	m.TotalTrackingHours = entity.TotalTrackingHours
	m.ProfitTargetHit = entity.ProfitTargetHit
	m.StopLossHit = entity.StopLossHit
	m.ClosedAt = entity.ClosedAt
}

// SignalKlineTrackingModel represents the signal_kline_tracking table
type SignalKlineTrackingModel struct {
	ID                    int64           `gorm:"column:id;primaryKey;autoIncrement"`
	SignalID              string          `gorm:"column:signal_id;size:36;not null;index:idx_signal_kline"`
	KlineOpenTime         time.Time       `gorm:"column:kline_open_time;not null;index:idx_signal_kline"`
	KlineCloseTime        time.Time       `gorm:"column:kline_close_time;not null"`
	HoursSinceSignal      decimal.Decimal `gorm:"column:hours_since_signal;type:decimal(10,2);not null"`
	OpenPrice             decimal.Decimal `gorm:"column:open_price;type:decimal(20,8);not null"`
	HighPrice             decimal.Decimal `gorm:"column:high_price;type:decimal(20,8);not null"`
	LowPrice              decimal.Decimal `gorm:"column:low_price;type:decimal(20,8);not null"`
	ClosePrice            decimal.Decimal `gorm:"column:close_price;type:decimal(20,8);not null"`
	Volume                decimal.Decimal `gorm:"column:volume;type:decimal(20,4);not null"`
	QuoteVolume           decimal.Decimal `gorm:"column:quote_volume;type:decimal(20,4);not null"`
	OpenChangePct         decimal.Decimal `gorm:"column:open_change_pct;type:decimal(10,4);not null"`
	HighChangePct         decimal.Decimal `gorm:"column:high_change_pct;type:decimal(10,4);not null"`
	LowChangePct          decimal.Decimal `gorm:"column:low_change_pct;type:decimal(10,4);not null"`
	CloseChangePct        decimal.Decimal `gorm:"column:close_change_pct;type:decimal(10,4);not null"`
	HourlyReturnPct       decimal.Decimal `gorm:"column:hourly_return_pct;type:decimal(10,4);not null"`
	MaxPotentialProfitPct decimal.Decimal `gorm:"column:max_potential_profit_pct;type:decimal(10,4);not null"`
	MaxPotentialLossPct   decimal.Decimal `gorm:"column:max_potential_loss_pct;type:decimal(10,4);not null"`
	IsProfitableAtHigh    bool            `gorm:"column:is_profitable_at_high;not null;default:false;index"`
	IsProfitableAtClose   bool            `gorm:"column:is_profitable_at_close;not null;default:false;index"`
	CreatedAt             time.Time       `gorm:"column:created_at;autoCreateTime"`
}

// TableName specifies the table name
func (SignalKlineTrackingModel) TableName() string {
	return "signal_kline_tracking"
}

// ToEntity converts model to domain entity
func (m *SignalKlineTrackingModel) ToEntity() *entity.SignalKlineTracking {
	return &entity.SignalKlineTracking{
		ID:                    m.ID,
		SignalID:              m.SignalID,
		KlineOpenTime:         m.KlineOpenTime,
		KlineCloseTime:        m.KlineCloseTime,
		HoursSinceSignal:      m.HoursSinceSignal,
		OpenPrice:             m.OpenPrice,
		HighPrice:             m.HighPrice,
		LowPrice:              m.LowPrice,
		ClosePrice:            m.ClosePrice,
		Volume:                m.Volume,
		QuoteVolume:           m.QuoteVolume,
		OpenChangePct:         m.OpenChangePct,
		HighChangePct:         m.HighChangePct,
		LowChangePct:          m.LowChangePct,
		CloseChangePct:        m.CloseChangePct,
		HourlyReturnPct:       m.HourlyReturnPct,
		MaxPotentialProfitPct: m.MaxPotentialProfitPct,
		MaxPotentialLossPct:   m.MaxPotentialLossPct,
		IsProfitableAtHigh:    m.IsProfitableAtHigh,
		IsProfitableAtClose:   m.IsProfitableAtClose,
		CreatedAt:             m.CreatedAt,
	}
}

// FromEntity converts domain entity to model
func (m *SignalKlineTrackingModel) FromEntity(entity *entity.SignalKlineTracking) {
	m.ID = entity.ID
	m.SignalID = entity.SignalID
	m.KlineOpenTime = entity.KlineOpenTime
	m.KlineCloseTime = entity.KlineCloseTime
	m.HoursSinceSignal = entity.HoursSinceSignal
	m.OpenPrice = entity.OpenPrice
	m.HighPrice = entity.HighPrice
	m.LowPrice = entity.LowPrice
	m.ClosePrice = entity.ClosePrice
	m.Volume = entity.Volume
	m.QuoteVolume = entity.QuoteVolume
	m.OpenChangePct = entity.OpenChangePct
	m.HighChangePct = entity.HighChangePct
	m.LowChangePct = entity.LowChangePct
	m.CloseChangePct = entity.CloseChangePct
	m.HourlyReturnPct = entity.HourlyReturnPct
	m.MaxPotentialProfitPct = entity.MaxPotentialProfitPct
	m.MaxPotentialLossPct = entity.MaxPotentialLossPct
	m.IsProfitableAtHigh = entity.IsProfitableAtHigh
	m.IsProfitableAtClose = entity.IsProfitableAtClose
}

// SignalRepository implements repository.SignalRepository
type SignalRepository struct {
	db *gorm.DB
}

// NewSignalRepository creates a new signal repository
func NewSignalRepository(db *gorm.DB) *SignalRepository {
	return &SignalRepository{db: db}
}

// Create creates a new signal
func (r *SignalRepository) Create(ctx context.Context, signal *entity.Signal) error {
	model := &SignalModel{}
	if err := model.FromEntity(signal); err != nil {
		return fmt.Errorf("failed to convert entity: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create signal: %w", err)
	}

	signal.ID = model.ID
	return nil
}

// Update updates an existing signal
func (r *SignalRepository) Update(ctx context.Context, signal *entity.Signal) error {
	model := &SignalModel{}
	if err := model.FromEntity(signal); err != nil {
		return fmt.Errorf("failed to convert entity: %w", err)
	}

	// Use Updates instead of Save to avoid updating zero-value CreatedAt
	// Updates will ignore zero values and only update specified fields
	if err := r.db.WithContext(ctx).Model(&SignalModel{}).
		Where("id = ?", model.ID).
		Updates(map[string]interface{}{
			"signal_id":           model.SignalID,
			"symbol":              model.Symbol,
			"signal_type":         model.Type,
			"strategy_name":       model.StrategyName,
			"generated_at":        model.GeneratedAt,
			"price_at_signal":     model.PriceAtSignal,
			"long_account_ratio":  model.LongAccountRatio,
			"short_account_ratio": model.ShortAccountRatio,
			"long_position_ratio": model.LongPositionRatio,
			"short_position_ratio": model.ShortPositionRatio,
			"long_trader_count":   model.LongTraderCount,
			"short_trader_count":  model.ShortTraderCount,
			"confirmation_start":  model.ConfirmationStart,
			"confirmation_end":    model.ConfirmationEnd,
			"is_confirmed":        model.IsConfirmed,
			"confirmed_at":        model.ConfirmedAt,
			"status":              model.Status,
			"reason":              model.Reason,
			"config_snapshot":     model.ConfigSnapshot,
		}).Error; err != nil {
		return fmt.Errorf("failed to update signal: %w", err)
	}

	return nil
}

// GetByID retrieves a signal by its UUID
func (r *SignalRepository) GetByID(ctx context.Context, signalID string) (*entity.Signal, error) {
	var model SignalModel
	if err := r.db.WithContext(ctx).Where("signal_id = ?", signalID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get signal: %w", err)
	}

	return model.ToEntity()
}

// GetBySymbol retrieves signals for a symbol
func (r *SignalRepository) GetBySymbol(ctx context.Context, symbol string, limit int) ([]*entity.Signal, error) {
	var models []SignalModel
	query := r.db.WithContext(ctx).Where("symbol = ?", symbol).Order("generated_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get signals by symbol: %w", err)
	}

	return r.modelsToEntities(models)
}

// GetByStatus retrieves signals with a specific status
func (r *SignalRepository) GetByStatus(ctx context.Context, status entity.SignalStatus, limit int) ([]*entity.Signal, error) {
	var models []SignalModel
	query := r.db.WithContext(ctx).Where("status = ?", string(status)).Order("generated_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get signals by status: %w", err)
	}

	return r.modelsToEntities(models)
}

// GetAll retrieves all signals
func (r *SignalRepository) GetAll(ctx context.Context) ([]*entity.Signal, error) {
	var models []SignalModel
	if err := r.db.WithContext(ctx).
		Order("generated_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all signals: %w", err)
	}

	return r.modelsToEntities(models)
}

// GetActiveSignals retrieves all active signals
func (r *SignalRepository) GetActiveSignals(ctx context.Context) ([]*entity.Signal, error) {
	var models []SignalModel
	if err := r.db.WithContext(ctx).
		Where("status IN ?", []string{
			string(entity.SignalStatusPending),
			string(entity.SignalStatusConfirmed),
			string(entity.SignalStatusTracking),
		}).
		Order("generated_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get active signals: %w", err)
	}

	return r.modelsToEntities(models)
}

// GetPendingSignals retrieves all pending signals
func (r *SignalRepository) GetPendingSignals(ctx context.Context) ([]*entity.Signal, error) {
	return r.GetByStatus(ctx, entity.SignalStatusPending, 0)
}

// GetConfirmedSignals retrieves all confirmed signals
func (r *SignalRepository) GetConfirmedSignals(ctx context.Context) ([]*entity.Signal, error) {
	return r.GetByStatus(ctx, entity.SignalStatusConfirmed, 0)
}

// GetTrackingSignals retrieves all signals being tracked
func (r *SignalRepository) GetTrackingSignals(ctx context.Context) ([]*entity.Signal, error) {
	return r.GetByStatus(ctx, entity.SignalStatusTracking, 0)
}

// GetRecentSignalsBySymbol retrieves recent signals for a symbol within a time window
func (r *SignalRepository) GetRecentSignalsBySymbol(ctx context.Context, symbol string, since time.Time) ([]*entity.Signal, error) {
	var models []SignalModel
	if err := r.db.WithContext(ctx).
		Where("symbol = ? AND generated_at >= ?", symbol, since).
		Order("generated_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent signals: %w", err)
	}

	return r.modelsToEntities(models)
}

// CountActiveSignalsBySymbol counts active signals for a symbol
func (r *SignalRepository) CountActiveSignalsBySymbol(ctx context.Context, symbol string) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&SignalModel{}).
		Where("symbol = ? AND status IN ?", symbol, []string{
			string(entity.SignalStatusPending),
			string(entity.SignalStatusConfirmed),
			string(entity.SignalStatusTracking),
		}).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count active signals: %w", err)
	}

	return int(count), nil
}

// GetSignalsInTimeRange retrieves signals generated within a time range
func (r *SignalRepository) GetSignalsInTimeRange(ctx context.Context, start, end time.Time) ([]*entity.Signal, error) {
	var models []SignalModel
	if err := r.db.WithContext(ctx).
		Where("generated_at >= ? AND generated_at <= ?", start, end).
		Order("generated_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get signals in time range: %w", err)
	}

	return r.modelsToEntities(models)
}

// GetSignalsByStrategy retrieves signals for a specific strategy
func (r *SignalRepository) GetSignalsByStrategy(ctx context.Context, strategyName string, limit int) ([]*entity.Signal, error) {
	var models []SignalModel
	query := r.db.WithContext(ctx).Where("strategy_name = ?", strategyName).Order("generated_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get signals by strategy: %w", err)
	}

	return r.modelsToEntities(models)
}

// CreateTracking creates a new signal tracking record
func (r *SignalRepository) CreateTracking(ctx context.Context, tracking *entity.SignalTracking) error {
	model := &SignalTrackingModel{}
	model.FromEntity(tracking)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create tracking: %w", err)
	}

	tracking.ID = model.ID
	return nil
}

// GetLatestTracking retrieves the latest tracking record for a signal
func (r *SignalRepository) GetLatestTracking(ctx context.Context, signalID string) (*entity.SignalTracking, error) {
	var model SignalTrackingModel
	if err := r.db.WithContext(ctx).
		Where("signal_id = ?", signalID).
		Order("tracked_at DESC").
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest tracking: %w", err)
	}

	return model.ToEntity(), nil
}

// GetAllTracking retrieves all tracking records for a signal
func (r *SignalRepository) GetAllTracking(ctx context.Context, signalID string) ([]*entity.SignalTracking, error) {
	var models []SignalTrackingModel
	if err := r.db.WithContext(ctx).
		Where("signal_id = ?", signalID).
		Order("tracked_at ASC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all tracking: %w", err)
	}

	trackings := make([]*entity.SignalTracking, len(models))
	for i, model := range models {
		trackings[i] = model.ToEntity()
	}

	return trackings, nil
}

// CreateOutcome creates a new signal outcome
func (r *SignalRepository) CreateOutcome(ctx context.Context, outcome *entity.SignalOutcome) error {
	model := &SignalOutcomeModel{}
	model.FromEntity(outcome)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create outcome: %w", err)
	}

	outcome.ID = model.ID
	return nil
}

// GetOutcome retrieves the outcome for a signal
func (r *SignalRepository) GetOutcome(ctx context.Context, signalID string) (*entity.SignalOutcome, error) {
	var model SignalOutcomeModel
	if err := r.db.WithContext(ctx).Where("signal_id = ?", signalID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get outcome: %w", err)
	}

	return model.ToEntity(), nil
}

// GetOutcomesBySignalIDs retrieves outcomes for multiple signals
func (r *SignalRepository) GetOutcomesBySignalIDs(ctx context.Context, signalIDs []string) (map[string]*entity.SignalOutcome, error) {
	if len(signalIDs) == 0 {
		return make(map[string]*entity.SignalOutcome), nil
	}

	var models []SignalOutcomeModel
	if err := r.db.WithContext(ctx).
		Where("signal_id IN ?", signalIDs).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get outcomes by signal IDs: %w", err)
	}

	outcomeMap := make(map[string]*entity.SignalOutcome, len(models))
	for _, model := range models {
		outcomeMap[model.SignalID] = model.ToEntity()
	}

	return outcomeMap, nil
}

// GetOutcomesByTimeRange retrieves outcomes within a time range
func (r *SignalRepository) GetOutcomesByTimeRange(ctx context.Context, start, end time.Time) ([]*entity.SignalOutcome, error) {
	var models []SignalOutcomeModel
	if err := r.db.WithContext(ctx).
		Where("closed_at >= ? AND closed_at <= ?", start, end).
		Order("closed_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get outcomes by time range: %w", err)
	}

	outcomes := make([]*entity.SignalOutcome, len(models))
	for i, model := range models {
		outcomes[i] = model.ToEntity()
	}

	return outcomes, nil
}

// GetOutcomesByStrategy retrieves outcomes for a specific strategy
func (r *SignalRepository) GetOutcomesByStrategy(ctx context.Context, strategyName string, start, end time.Time) ([]*entity.SignalOutcome, error) {
	var models []SignalOutcomeModel
	if err := r.db.WithContext(ctx).
		Table("signal_outcomes").
		Joins("INNER JOIN signals ON signal_outcomes.signal_id = signals.signal_id").
		Where("signals.strategy_name = ? AND signal_outcomes.closed_at >= ? AND signal_outcomes.closed_at <= ?", strategyName, start, end).
		Order("signal_outcomes.closed_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get outcomes by strategy: %w", err)
	}

	outcomes := make([]*entity.SignalOutcome, len(models))
	for i, model := range models {
		outcomes[i] = model.ToEntity()
	}

	return outcomes, nil
}

// modelsToEntities converts signal models to entities
func (r *SignalRepository) modelsToEntities(models []SignalModel) ([]*entity.Signal, error) {
	signals := make([]*entity.Signal, len(models))
	for i, model := range models {
		signal, err := model.ToEntity()
		if err != nil {
			return nil, err
		}
		signals[i] = signal
	}
	return signals, nil
}

// Kline tracking methods

// CreateKlineTracking creates a new kline tracking record
func (r *SignalRepository) CreateKlineTracking(ctx context.Context, tracking *entity.SignalKlineTracking) error {
	model := &SignalKlineTrackingModel{}
	model.FromEntity(tracking)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create kline tracking: %w", err)
	}

	tracking.ID = model.ID
	return nil
}

// GetLatestKlineTracking retrieves the latest kline tracking record for a signal
func (r *SignalRepository) GetLatestKlineTracking(ctx context.Context, signalID string) (*entity.SignalKlineTracking, error) {
	var model SignalKlineTrackingModel
	if err := r.db.WithContext(ctx).
		Where("signal_id = ?", signalID).
		Order("kline_open_time DESC").
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest kline tracking: %w", err)
	}

	return model.ToEntity(), nil
}

// GetKlineTrackingBySignal retrieves all kline tracking records for a signal
func (r *SignalRepository) GetKlineTrackingBySignal(ctx context.Context, signalID string) ([]*entity.SignalKlineTracking, error) {
	var models []SignalKlineTrackingModel
	if err := r.db.WithContext(ctx).
		Where("signal_id = ?", signalID).
		Order("kline_open_time ASC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get kline tracking by signal: %w", err)
	}

	trackings := make([]*entity.SignalKlineTracking, len(models))
	for i, model := range models {
		trackings[i] = model.ToEntity()
	}

	return trackings, nil
}

// GetKlineTrackingInTimeRange retrieves kline tracking records within a time range
func (r *SignalRepository) GetKlineTrackingInTimeRange(ctx context.Context, start, end time.Time) ([]*entity.SignalKlineTracking, error) {
	var models []SignalKlineTrackingModel
	if err := r.db.WithContext(ctx).
		Where("kline_open_time >= ? AND kline_open_time <= ?", start, end).
		Order("kline_open_time ASC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get kline tracking in time range: %w", err)
	}

	trackings := make([]*entity.SignalKlineTracking, len(models))
	for i, model := range models {
		trackings[i] = model.ToEntity()
	}

	return trackings, nil
}
