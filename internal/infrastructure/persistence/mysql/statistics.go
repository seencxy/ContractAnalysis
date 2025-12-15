package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"ContractAnalysis/internal/domain/repository"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StrategyStatisticsModel represents the strategy_statistics table
type StrategyStatisticsModel struct {
	ID                 int64            `gorm:"column:id;primaryKey;autoIncrement"`
	StrategyName       string           `gorm:"column:strategy_name;size:50;not null;index"`
	Symbol             sql.NullString   `gorm:"column:symbol;size:50;index"`
	PeriodStart        time.Time        `gorm:"column:period_start;not null;index:idx_period_range"`
	PeriodEnd          time.Time        `gorm:"column:period_end;not null;index:idx_period_range"`
	PeriodLabel        string           `gorm:"column:period_label;size:20;not null;index"`
	TotalSignals       int              `gorm:"column:total_signals;default:0"`
	ConfirmedSignals   int              `gorm:"column:confirmed_signals;default:0"`
	InvalidatedSignals int              `gorm:"column:invalidated_signals;default:0"`
	ProfitableSignals  int              `gorm:"column:profitable_signals;default:0"`
	LosingSignals      int              `gorm:"column:losing_signals;default:0"`
	NeutralSignals     int              `gorm:"column:neutral_signals;default:0"`
	WinRate            *decimal.Decimal `gorm:"column:win_rate;type:decimal(10,4)"`
	AvgProfitPct       *decimal.Decimal `gorm:"column:avg_profit_pct;type:decimal(10,4)"`
	AvgLossPct         *decimal.Decimal `gorm:"column:avg_loss_pct;type:decimal(10,4)"`
	AvgHoldingHours    *decimal.Decimal `gorm:"column:avg_holding_hours;type:decimal(10,2)"`
	BestSignalPct      *decimal.Decimal `gorm:"column:best_signal_pct;type:decimal(10,4)"`
	WorstSignalPct     *decimal.Decimal `gorm:"column:worst_signal_pct;type:decimal(10,4)"`
	ProfitFactor       *decimal.Decimal `gorm:"column:profit_factor;type:decimal(10,4)"`

	// Kline-based win rate metrics
	KlineTheoreticalWinRate   *decimal.Decimal `gorm:"column:kline_theoretical_win_rate;type:decimal(10,4)"`
	KlineCloseWinRate         *decimal.Decimal `gorm:"column:kline_close_win_rate;type:decimal(10,4)"`
	TotalKlineHours           int              `gorm:"column:total_kline_hours;default:0"`
	ProfitableKlineHoursHigh  int              `gorm:"column:profitable_kline_hours_high;default:0"`
	ProfitableKlineHoursClose int              `gorm:"column:profitable_kline_hours_close;default:0"`

	// Hourly return statistics
	AvgHourlyReturnPct *decimal.Decimal `gorm:"column:avg_hourly_return_pct;type:decimal(10,4)"`
	MaxHourlyReturnPct *decimal.Decimal `gorm:"column:max_hourly_return_pct;type:decimal(10,4)"`
	MinHourlyReturnPct *decimal.Decimal `gorm:"column:min_hourly_return_pct;type:decimal(10,4)"`

	// Theoretical maximum profit/loss
	AvgMaxPotentialProfitPct *decimal.Decimal `gorm:"column:avg_max_potential_profit_pct;type:decimal(10,4)"`
	AvgMaxPotentialLossPct   *decimal.Decimal `gorm:"column:avg_max_potential_loss_pct;type:decimal(10,4)"`

	CalculatedAt time.Time `gorm:"column:calculated_at;autoCreateTime;index"`
}

// TableName specifies the table name
func (StrategyStatisticsModel) TableName() string {
	return "strategy_statistics"
}

// ToEntity converts model to domain entity
func (m *StrategyStatisticsModel) ToEntity() *repository.StrategyStatistics {
	var symbol *string
	if m.Symbol.Valid {
		symbol = &m.Symbol.String
	}

	return &repository.StrategyStatistics{
		ID:                 m.ID,
		StrategyName:       m.StrategyName,
		Symbol:             symbol,
		PeriodStart:        m.PeriodStart,
		PeriodEnd:          m.PeriodEnd,
		PeriodLabel:        m.PeriodLabel,
		TotalSignals:       m.TotalSignals,
		ConfirmedSignals:   m.ConfirmedSignals,
		InvalidatedSignals: m.InvalidatedSignals,
		ProfitableSignals:  m.ProfitableSignals,
		LosingSignals:      m.LosingSignals,
		NeutralSignals:     m.NeutralSignals,
		WinRate:            m.WinRate,
		AvgProfitPct:       m.AvgProfitPct,
		AvgLossPct:         m.AvgLossPct,
		AvgHoldingHours:    m.AvgHoldingHours,
		BestSignalPct:      m.BestSignalPct,
		WorstSignalPct:     m.WorstSignalPct,
		ProfitFactor:       m.ProfitFactor,

		// Kline-based win rate metrics
		KlineTheoreticalWinRate:   m.KlineTheoreticalWinRate,
		KlineCloseWinRate:         m.KlineCloseWinRate,
		TotalKlineHours:           m.TotalKlineHours,
		ProfitableKlineHoursHigh:  m.ProfitableKlineHoursHigh,
		ProfitableKlineHoursClose: m.ProfitableKlineHoursClose,

		// Hourly return statistics
		AvgHourlyReturnPct: m.AvgHourlyReturnPct,
		MaxHourlyReturnPct: m.MaxHourlyReturnPct,
		MinHourlyReturnPct: m.MinHourlyReturnPct,

		// Theoretical maximum profit/loss
		AvgMaxPotentialProfitPct: m.AvgMaxPotentialProfitPct,
		AvgMaxPotentialLossPct:   m.AvgMaxPotentialLossPct,

		CalculatedAt: m.CalculatedAt,
	}
}

// FromEntity converts domain entity to model
func (m *StrategyStatisticsModel) FromEntity(entity *repository.StrategyStatistics) {
	m.ID = entity.ID
	m.StrategyName = entity.StrategyName

	if entity.Symbol != nil {
		m.Symbol = sql.NullString{String: *entity.Symbol, Valid: true}
	} else {
		m.Symbol = sql.NullString{Valid: false}
	}

	m.PeriodStart = entity.PeriodStart
	m.PeriodEnd = entity.PeriodEnd
	m.PeriodLabel = entity.PeriodLabel
	m.TotalSignals = entity.TotalSignals
	m.ConfirmedSignals = entity.ConfirmedSignals
	m.InvalidatedSignals = entity.InvalidatedSignals
	m.ProfitableSignals = entity.ProfitableSignals
	m.LosingSignals = entity.LosingSignals
	m.NeutralSignals = entity.NeutralSignals
	m.WinRate = entity.WinRate
	m.AvgProfitPct = entity.AvgProfitPct
	m.AvgLossPct = entity.AvgLossPct
	m.AvgHoldingHours = entity.AvgHoldingHours
	m.BestSignalPct = entity.BestSignalPct
	m.WorstSignalPct = entity.WorstSignalPct
	m.ProfitFactor = entity.ProfitFactor

	// Kline-based win rate metrics
	m.KlineTheoreticalWinRate = entity.KlineTheoreticalWinRate
	m.KlineCloseWinRate = entity.KlineCloseWinRate
	m.TotalKlineHours = entity.TotalKlineHours
	m.ProfitableKlineHoursHigh = entity.ProfitableKlineHoursHigh
	m.ProfitableKlineHoursClose = entity.ProfitableKlineHoursClose

	// Hourly return statistics
	m.AvgHourlyReturnPct = entity.AvgHourlyReturnPct
	m.MaxHourlyReturnPct = entity.MaxHourlyReturnPct
	m.MinHourlyReturnPct = entity.MinHourlyReturnPct

	// Theoretical maximum profit/loss
	m.AvgMaxPotentialProfitPct = entity.AvgMaxPotentialProfitPct
	m.AvgMaxPotentialLossPct = entity.AvgMaxPotentialLossPct
}

// StatisticsRepository implements repository.StatisticsRepository
type StatisticsRepository struct {
	db *gorm.DB
}

// NewStatisticsRepository creates a new statistics repository
func NewStatisticsRepository(db *gorm.DB) repository.StatisticsRepository {
	return &StatisticsRepository{db: db}
}

// Create creates a new statistics record
func (r *StatisticsRepository) Create(ctx context.Context, stats *repository.StrategyStatistics) error {
	model := &StrategyStatisticsModel{}
	model.FromEntity(stats)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create statistics: %w", err)
	}

	stats.ID = model.ID
	return nil
}

// CreateOrUpdate creates or updates a statistics record
func (r *StatisticsRepository) CreateOrUpdate(ctx context.Context, stats *repository.StrategyStatistics) error {
	model := &StrategyStatisticsModel{}
	model.FromEntity(stats)

	// Use UPSERT to create or update based on unique constraint
	symbolValue := ""
	if stats.Symbol != nil {
		symbolValue = *stats.Symbol
	}

	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "strategy_name"},
				{Name: "symbol"},
				{Name: "period_label"},
				{Name: "period_start"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"period_end",
				"total_signals",
				"confirmed_signals",
				"invalidated_signals",
				"profitable_signals",
				"losing_signals",
				"neutral_signals",
				"win_rate",
				"avg_profit_pct",
				"avg_loss_pct",
				"avg_holding_hours",
				"best_signal_pct",
				"worst_signal_pct",
				"profit_factor",
				"kline_theoretical_win_rate",
				"kline_close_win_rate",
				"total_kline_hours",
				"profitable_kline_hours_high",
				"profitable_kline_hours_close",
				"avg_hourly_return_pct",
				"max_hourly_return_pct",
				"min_hourly_return_pct",
				"avg_max_potential_profit_pct",
				"avg_max_potential_loss_pct",
				"calculated_at",
			}),
		}).
		Where("strategy_name = ? AND COALESCE(symbol, '') = ? AND period_label = ? AND period_start = ?",
			stats.StrategyName, symbolValue, stats.PeriodLabel, stats.PeriodStart).
		Create(model).Error; err != nil {
		return fmt.Errorf("failed to create or update statistics: %w", err)
	}

	stats.ID = model.ID
	return nil
}

// GetByStrategyAndPeriod retrieves statistics for a strategy and period
func (r *StatisticsRepository) GetByStrategyAndPeriod(ctx context.Context, strategyName, periodLabel string, symbol *string) (*repository.StrategyStatistics, error) {
	var model StrategyStatisticsModel

	query := r.db.WithContext(ctx).
		Where("strategy_name = ? AND period_label = ?", strategyName, periodLabel).
		Order("calculated_at DESC")

	if symbol != nil {
		query = query.Where("symbol = ?", *symbol)
	} else {
		query = query.Where("symbol IS NULL")
	}

	if err := query.First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByStrategy retrieves all statistics for a strategy
func (r *StatisticsRepository) GetByStrategy(ctx context.Context, strategyName string) ([]*repository.StrategyStatistics, error) {
	var models []StrategyStatisticsModel
	if err := r.db.WithContext(ctx).
		Where("strategy_name = ?", strategyName).
		Order("period_label, calculated_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get statistics by strategy: %w", err)
	}

	stats := make([]*repository.StrategyStatistics, len(models))
	for i, model := range models {
		stats[i] = model.ToEntity()
	}

	return stats, nil
}

// GetByPeriod retrieves all statistics for a period
func (r *StatisticsRepository) GetByPeriod(ctx context.Context, periodLabel string) ([]*repository.StrategyStatistics, error) {
	var models []StrategyStatisticsModel
	if err := r.db.WithContext(ctx).
		Where("period_label = ?", periodLabel).
		Order("strategy_name, calculated_at DESC").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get statistics by period: %w", err)
	}

	stats := make([]*repository.StrategyStatistics, len(models))
	for i, model := range models {
		stats[i] = model.ToEntity()
	}

	return stats, nil
}

// GetLatest retrieves the latest statistics for each strategy and period
func (r *StatisticsRepository) GetLatest(ctx context.Context) ([]*repository.StrategyStatistics, error) {
	// Get latest statistics grouped by strategy_name, symbol, and period_label
	var models []StrategyStatisticsModel

	subQuery := r.db.Model(&StrategyStatisticsModel{}).
		Select("strategy_name, COALESCE(symbol, '') as symbol, period_label, MAX(calculated_at) as max_calc").
		Group("strategy_name, COALESCE(symbol, ''), period_label")

	if err := r.db.WithContext(ctx).
		Table("strategy_statistics as ss").
		Joins("INNER JOIN (?) as latest ON ss.strategy_name = latest.strategy_name AND COALESCE(ss.symbol, '') = latest.symbol AND ss.period_label = latest.period_label AND ss.calculated_at = latest.max_calc", subQuery).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get latest statistics: %w", err)
	}

	stats := make([]*repository.StrategyStatistics, len(models))
	for i, model := range models {
		stats[i] = model.ToEntity()
	}

	return stats, nil
}

// GetPreviousCalculation retrieves the statistics calculation before the given one
func (r *StatisticsRepository) GetPreviousCalculation(
	ctx context.Context,
	strategyName, periodLabel string,
	symbol *string,
	currentCalculatedAt time.Time,
) (*repository.StrategyStatistics, error) {
	var model StrategyStatisticsModel

	query := r.db.WithContext(ctx).
		Where("strategy_name = ?", strategyName).
		Where("period_label = ?", periodLabel).
		Where("calculated_at < ?", currentCalculatedAt).
		Order("calculated_at DESC").
		Limit(1)

	if symbol != nil {
		query = query.Where("symbol = ?", *symbol)
	} else {
		query = query.Where("symbol IS NULL")
	}

	if err := query.First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No previous record
		}
		return nil, fmt.Errorf("failed to get previous calculation: %w", err)
	}

	return model.ToEntity(), nil
}

// GetByTimeRange retrieves statistics within a time range
func (r *StatisticsRepository) GetByTimeRange(
	ctx context.Context,
	startTime, endTime time.Time,
	strategyName, symbol *string,
) ([]*repository.StrategyStatistics, error) {
	var models []StrategyStatisticsModel

	query := r.db.WithContext(ctx).
		Where("calculated_at >= ?", startTime).
		Where("calculated_at <= ?", endTime).
		Order("calculated_at DESC")

	// Optional strategy filter
	if strategyName != nil && *strategyName != "" {
		query = query.Where("strategy_name = ?", *strategyName)
	}

	// Optional symbol filter
	if symbol != nil && *symbol != "" {
		query = query.Where("symbol = ?", *symbol)
	} else if symbol != nil && *symbol == "" {
		// Empty string means filter for overall stats (NULL symbol)
		query = query.Where("symbol IS NULL")
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get statistics by time range: %w", err)
	}

	stats := make([]*repository.StrategyStatistics, len(models))
	for i, model := range models {
		stats[i] = model.ToEntity()
	}

	return stats, nil
}

// DeleteOlderThan deletes statistics older than the specified time
func (r *StatisticsRepository) DeleteOlderThan(ctx context.Context, before time.Time) error {
	if err := r.db.WithContext(ctx).
		Where("calculated_at < ?", before).
		Delete(&StrategyStatisticsModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete old statistics: %w", err)
	}

	return nil
}
