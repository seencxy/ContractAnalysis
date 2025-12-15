package mysql

import (
	"context"
	"fmt"
	"time"

	"ContractAnalysis/internal/domain/repository"

	"gorm.io/gorm"
)

// TradingPairModel represents the trading_pairs table
type TradingPairModel struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Symbol    string    `gorm:"column:symbol;uniqueIndex;size:50;not null"`
	IsActive  bool      `gorm:"column:is_active;default:true"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name
func (TradingPairModel) TableName() string {
	return "trading_pairs"
}

// ToEntity converts model to domain entity
func (m *TradingPairModel) ToEntity() *repository.TradingPair {
	return &repository.TradingPair{
		ID:       m.ID,
		Symbol:   m.Symbol,
		IsActive: m.IsActive,
	}
}

// FromEntity converts domain entity to model
func (m *TradingPairModel) FromEntity(entity *repository.TradingPair) {
	m.ID = entity.ID
	m.Symbol = entity.Symbol
	m.IsActive = entity.IsActive
}

// TradingPairRepository implements repository.TradingPairRepository
type TradingPairRepository struct {
	db *gorm.DB
}

// NewTradingPairRepository creates a new trading pair repository
func NewTradingPairRepository(db *gorm.DB) repository.TradingPairRepository {
	return &TradingPairRepository{db: db}
}

// Create creates a new trading pair
func (r *TradingPairRepository) Create(ctx context.Context, pair *repository.TradingPair) error {
	model := &TradingPairModel{}
	model.FromEntity(pair)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create trading pair: %w", err)
	}

	pair.ID = model.ID
	return nil
}

// CreateBatch creates multiple trading pairs in a batch
func (r *TradingPairRepository) CreateBatch(ctx context.Context, pairs []*repository.TradingPair) error {
	if len(pairs) == 0 {
		return nil
	}

	models := make([]*TradingPairModel, len(pairs))
	for i, pair := range pairs {
		model := &TradingPairModel{}
		model.FromEntity(pair)
		models[i] = model
	}

	if err := r.db.WithContext(ctx).Create(&models).Error; err != nil {
		return fmt.Errorf("failed to create trading pairs batch: %w", err)
	}

	// Update IDs
	for i, model := range models {
		pairs[i].ID = model.ID
	}

	return nil
}

// Update updates an existing trading pair
func (r *TradingPairRepository) Update(ctx context.Context, pair *repository.TradingPair) error {
	model := &TradingPairModel{}
	model.FromEntity(pair)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update trading pair: %w", err)
	}

	return nil
}

// GetBySymbol retrieves a trading pair by symbol
func (r *TradingPairRepository) GetBySymbol(ctx context.Context, symbol string) (*repository.TradingPair, error) {
	var model TradingPairModel
	if err := r.db.WithContext(ctx).Where("symbol = ?", symbol).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get trading pair: %w", err)
	}

	return model.ToEntity(), nil
}

// GetAll retrieves all trading pairs
func (r *TradingPairRepository) GetAll(ctx context.Context) ([]*repository.TradingPair, error) {
	var models []TradingPairModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all trading pairs: %w", err)
	}

	pairs := make([]*repository.TradingPair, len(models))
	for i, model := range models {
		pairs[i] = model.ToEntity()
	}

	return pairs, nil
}

// GetActive retrieves all active trading pairs
func (r *TradingPairRepository) GetActive(ctx context.Context) ([]*repository.TradingPair, error) {
	var models []TradingPairModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get active trading pairs: %w", err)
	}

	pairs := make([]*repository.TradingPair, len(models))
	for i, model := range models {
		pairs[i] = model.ToEntity()
	}

	return pairs, nil
}

// SetActive sets the active status of a trading pair
func (r *TradingPairRepository) SetActive(ctx context.Context, symbol string, isActive bool) error {
	if err := r.db.WithContext(ctx).
		Model(&TradingPairModel{}).
		Where("symbol = ?", symbol).
		Update("is_active", isActive).Error; err != nil {
		return fmt.Errorf("failed to set active status: %w", err)
	}

	return nil
}

// Exists checks if a trading pair exists
func (r *TradingPairRepository) Exists(ctx context.Context, symbol string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&TradingPairModel{}).
		Where("symbol = ?", symbol).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}

	return count > 0, nil
}
