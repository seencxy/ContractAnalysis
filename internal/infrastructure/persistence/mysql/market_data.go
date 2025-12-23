package mysql

import (
	"context"
	"fmt"
	"time"

	"ContractAnalysis/internal/domain/entity"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// MarketDataModel represents the market_data table
type MarketDataModel struct {
	ID                     int64           `gorm:"column:id;primaryKey;autoIncrement"`
	Symbol                 string          `gorm:"column:symbol;size:50;not null;index:idx_symbol_timestamp"`
	Timestamp              time.Time       `gorm:"column:timestamp;not null;uniqueIndex:uk_symbol_timestamp;index:idx_symbol_timestamp"`
	LongAccountRatio       decimal.Decimal `gorm:"column:long_account_ratio;type:decimal(10,4);not null"`
	ShortAccountRatio      decimal.Decimal `gorm:"column:short_account_ratio;type:decimal(10,4);not null"`
	LongPositionRatio      decimal.Decimal `gorm:"column:long_position_ratio;type:decimal(10,4);not null"`
	ShortPositionRatio     decimal.Decimal `gorm:"column:short_position_ratio;type:decimal(10,4);not null"`
	PositionRatioAvailable bool            `gorm:"column:position_ratio_available;default:true"`
	DataQualityScore       int             `gorm:"column:data_quality_score;type:tinyint;default:100"`
	Price                  decimal.Decimal `gorm:"column:price;type:decimal(20,8);not null"`
	Volume24h              decimal.Decimal `gorm:"column:volume_24h;type:decimal(20,2)"`
	OpenInterest           decimal.Decimal `gorm:"column:open_interest;type:decimal(20,8);default:0"`
	FundingRate            decimal.Decimal `gorm:"column:funding_rate;type:decimal(10,8);default:0"`
	CreatedAt              time.Time       `gorm:"column:created_at;autoCreateTime"`
}

// TableName specifies the table name
func (MarketDataModel) TableName() string {
	return "market_data"
}

// ToEntity converts model to domain entity
func (m *MarketDataModel) ToEntity() *entity.MarketData {
	return &entity.MarketData{
		ID:                     m.ID,
		Symbol:                 m.Symbol,
		Timestamp:              m.Timestamp,
		LongAccountRatio:       m.LongAccountRatio,
		ShortAccountRatio:      m.ShortAccountRatio,
		LongPositionRatio:      m.LongPositionRatio,
		ShortPositionRatio:     m.ShortPositionRatio,
		PositionRatioAvailable: m.PositionRatioAvailable,
		DataQualityScore:       m.DataQualityScore,
		Price:                  m.Price,
		Volume24h:              m.Volume24h,
		OpenInterest:           m.OpenInterest,
		FundingRate:            m.FundingRate,
		CreatedAt:              m.CreatedAt,
	}
}

// FromEntity converts domain entity to model
func (m *MarketDataModel) FromEntity(entity *entity.MarketData) {
	m.ID = entity.ID
	m.Symbol = entity.Symbol
	m.Timestamp = entity.Timestamp
	m.LongAccountRatio = entity.LongAccountRatio
	m.ShortAccountRatio = entity.ShortAccountRatio
	m.LongPositionRatio = entity.LongPositionRatio
	m.ShortPositionRatio = entity.ShortPositionRatio
	m.PositionRatioAvailable = entity.PositionRatioAvailable
	m.DataQualityScore = entity.DataQualityScore
	m.Price = entity.Price
	m.Volume24h = entity.Volume24h
	m.OpenInterest = entity.OpenInterest
	m.FundingRate = entity.FundingRate
}

// MarketDataRepository implements repository.MarketDataRepository
type MarketDataRepository struct {
	db *gorm.DB
}

// NewMarketDataRepository creates a new market data repository
func NewMarketDataRepository(db *gorm.DB) *MarketDataRepository {
	return &MarketDataRepository{db: db}
}

// Create creates a new market data record
func (r *MarketDataRepository) Create(ctx context.Context, data *entity.MarketData) error {
	model := &MarketDataModel{}
	model.FromEntity(data)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create market data: %w", err)
	}

	data.ID = model.ID
	return nil
}

// CreateBatch creates multiple market data records in a batch
func (r *MarketDataRepository) CreateBatch(ctx context.Context, dataList []*entity.MarketData) error {
	if len(dataList) == 0 {
		return nil
	}

	models := make([]*MarketDataModel, len(dataList))
	for i, data := range dataList {
		model := &MarketDataModel{}
		model.FromEntity(data)
		models[i] = model
	}

	// Use batch insert for better performance
	batchSize := 100
	if err := r.db.WithContext(ctx).CreateInBatches(models, batchSize).Error; err != nil {
		return fmt.Errorf("failed to create market data batch: %w", err)
	}

	// Update IDs
	for i, model := range models {
		dataList[i].ID = model.ID
	}

	return nil
}

// GetBySymbol retrieves market data for a symbol within a time range
func (r *MarketDataRepository) GetBySymbol(ctx context.Context, symbol string, start, end time.Time) ([]*entity.MarketData, error) {
	var models []MarketDataModel
	if err := r.db.WithContext(ctx).
		Where("symbol = ? AND timestamp >= ? AND timestamp <= ?", symbol, start, end).
		Order("timestamp DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get market data by symbol: %w", err)
	}

	dataList := make([]*entity.MarketData, len(models))
	for i, model := range models {
		dataList[i] = model.ToEntity()
	}

	return dataList, nil
}

// GetLatestBySymbol retrieves the latest market data for a symbol
func (r *MarketDataRepository) GetLatestBySymbol(ctx context.Context, symbol string) (*entity.MarketData, error) {
	var model MarketDataModel
	if err := r.db.WithContext(ctx).
		Where("symbol = ?", symbol).
		Order("timestamp DESC").
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest market data: %w", err)
	}

	return model.ToEntity(), nil
}

// GetLatestForAllSymbols retrieves the latest market data for all symbols
func (r *MarketDataRepository) GetLatestForAllSymbols(ctx context.Context) ([]*entity.MarketData, error) {
	// Use subquery to get the latest timestamp for each symbol
	var models []MarketDataModel
	subQuery := r.db.Model(&MarketDataModel{}).
		Select("symbol, MAX(timestamp) as max_timestamp").
		Group("symbol")

	if err := r.db.WithContext(ctx).
		Table("market_data as md").
		Joins("INNER JOIN (?) as latest ON md.symbol = latest.symbol AND md.timestamp = latest.max_timestamp", subQuery).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get latest market data for all symbols: %w", err)
	}

	dataList := make([]*entity.MarketData, len(models))
	for i, model := range models {
		dataList[i] = model.ToEntity()
	}

	return dataList, nil
}

// GetRecentBySymbol retrieves the most recent N records for a symbol
func (r *MarketDataRepository) GetRecentBySymbol(ctx context.Context, symbol string, limit int) ([]*entity.MarketData, error) {
	var models []MarketDataModel
	if err := r.db.WithContext(ctx).
		Where("symbol = ?", symbol).
		Order("timestamp DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent market data: %w", err)
	}

	dataList := make([]*entity.MarketData, len(models))
	for i, model := range models {
		dataList[i] = model.ToEntity()
	}

	return dataList, nil
}

// DeleteOlderThan deletes market data older than the specified time
func (r *MarketDataRepository) DeleteOlderThan(ctx context.Context, before time.Time) error {
	if err := r.db.WithContext(ctx).
		Where("timestamp < ?", before).
		Delete(&MarketDataModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete old market data: %w", err)
	}

	return nil
}

// Count returns the total number of market data records
func (r *MarketDataRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&MarketDataModel{}).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count market data: %w", err)
	}

	return count, nil
}
