package usecase

import (
	"context"
	"fmt"
	"time"

	"ContractAnalysis/config"
	"ContractAnalysis/internal/domain/entity"
	"ContractAnalysis/internal/domain/repository"
	"ContractAnalysis/internal/infrastructure/binance"
	"ContractAnalysis/internal/infrastructure/logger"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// Collector orchestrates the data collection process
type Collector struct {
	binanceClient   *binance.Client
	marketDataRepo  *repository.MarketDataRepository
	tradingPairRepo repository.TradingPairRepository
	config          config.CollectionConfig
	logger          *logger.Logger
}

// NewCollector creates a new collector
func NewCollector(
	binanceClient *binance.Client,
	marketDataRepo *repository.MarketDataRepository,
	tradingPairRepo repository.TradingPairRepository,
	cfg config.CollectionConfig,
) *Collector {
	return &Collector{
		binanceClient:   binanceClient,
		marketDataRepo:  marketDataRepo,
		tradingPairRepo: tradingPairRepo,
		config:          cfg,
		logger:          logger.WithComponent("collector"),
	}
}

// CollectAll collects market data for all active trading pairs
func (c *Collector) CollectAll(ctx context.Context) error {
	if !c.config.Enabled {
		c.logger.Info("Data collection is disabled")
		return nil
	}

	c.logger.Info("Starting data collection")
	startTime := time.Now()

	// Get all USDT futures pairs from Binance
	allPairs, err := c.binanceClient.GetAllUSDTFuturesPairs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get trading pairs: %w", err)
	}

	// Filter excluded pairs
	pairs := c.filterPairs(allPairs)
	c.logger.Info("Found trading pairs",
		zap.Int("total", len(allPairs)),
		zap.Int("filtered", len(pairs)),
	)

	// Update trading pairs in database
	if err := c.updateTradingPairs(ctx, pairs); err != nil {
		c.logger.WithError(err).Warn("Failed to update trading pairs")
		// Continue even if this fails
	}

	// Collect market data for each pair
	collected := 0
	failed := 0
	failedSymbols := make([]string, 0)

	for _, symbol := range pairs {
		if err := c.collectForSymbol(ctx, symbol); err != nil {
			c.logger.WithError(err).WithSymbol(symbol).Warn("Failed to collect data for symbol")
			failed++
			failedSymbols = append(failedSymbols, symbol)
			continue
		}
		collected++

		// Small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	duration := time.Since(startTime)
	totalPairs := len(pairs)
	successRate := float64(collected) / float64(totalPairs) * 100

	c.logger.Info("Data collection completed",
		zap.Int("total_pairs", totalPairs),
		zap.Int("collected", collected),
		zap.Int("failed", failed),
		zap.Float64("success_rate", successRate),
		zap.Duration("duration", duration),
		zap.Strings("failed_symbols", failedSymbols),
	)

	// Warning if success rate is low
	if successRate < 95.0 {
		c.logger.Warn("Low data collection success rate detected",
			zap.Float64("success_rate", successRate),
			zap.Int("failed_count", failed),
		)
	}

	// Error if success rate is critically low
	if successRate < 80.0 {
		c.logger.Error("Critically low data collection success rate",
			zap.Float64("success_rate", successRate),
			zap.Int("total", totalPairs),
			zap.Int("collected", collected),
			zap.Int("failed", failed),
		)
	}

	if failed > 0 && collected == 0 {
		return fmt.Errorf("failed to collect data for all symbols")
	}

	return nil
}

// CollectForSymbol collects market data for a specific symbol
func (c *Collector) CollectForSymbol(ctx context.Context, symbol string) error {
	return c.collectForSymbol(ctx, symbol)
}

// collectForSymbol collects and stores market data for a symbol
func (c *Collector) collectForSymbol(ctx context.Context, symbol string) error {
	c.logger.Debug("Collecting data for symbol", zap.String("symbol", symbol))

	// Fetch market data from Binance with retry
	var marketData *binance.MarketData
	var err error

	for attempt := 0; attempt < c.config.Retry.MaxAttempts; attempt++ {
		marketData, err = c.binanceClient.GetMarketData(ctx, symbol)
		if err == nil {
			break
		}

		if attempt < c.config.Retry.MaxAttempts-1 {
			delay := c.config.Retry.Delay * time.Duration(attempt+1)
			c.logger.Debug("Retrying after error",
				zap.String("symbol", symbol),
				zap.Int("attempt", attempt+1),
				zap.String("delay", delay.String()),
			)
			time.Sleep(delay)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to fetch market data after %d attempts: %w", c.config.Retry.MaxAttempts, err)
	}

	// Convert to domain entity
	entity := c.convertToEntity(marketData)

	// Validate entity
	if err := entity.Validate(); err != nil {
		return fmt.Errorf("invalid market data: %w", err)
	}

	// Store in database
	repo := *c.marketDataRepo
	if err := repo.Create(ctx, entity); err != nil {
		return fmt.Errorf("failed to store market data: %w", err)
	}

	c.logger.Debug("Successfully collected data for symbol",
		zap.String("symbol", symbol),
	)

	return nil
}

// updateTradingPairs updates the trading pairs in the database
func (c *Collector) updateTradingPairs(ctx context.Context, symbols []string) error {
	// Get existing pairs
	existingPairs, err := c.tradingPairRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get existing pairs: %w", err)
	}

	existingMap := make(map[string]bool)
	for _, pair := range existingPairs {
		existingMap[pair.Symbol] = true
	}

	// Create new pairs
	var newPairs []*repository.TradingPair
	for _, symbol := range symbols {
		if !existingMap[symbol] {
			newPairs = append(newPairs, &repository.TradingPair{
				Symbol:   symbol,
				IsActive: true,
			})
		}
	}

	if len(newPairs) > 0 {
		if err := c.tradingPairRepo.CreateBatch(ctx, newPairs); err != nil {
			return fmt.Errorf("failed to create new pairs: %w", err)
		}
		c.logger.Info("Created new trading pairs", zap.Int("count", len(newPairs)))
	}

	return nil
}

// filterPairs filters out excluded pairs
func (c *Collector) filterPairs(pairs []string) []string {
	if len(c.config.PairFilter.ExcludePairs) == 0 {
		return pairs
	}

	excludeMap := make(map[string]bool)
	for _, symbol := range c.config.PairFilter.ExcludePairs {
		excludeMap[symbol] = true
	}

	var filtered []string
	for _, symbol := range pairs {
		if !excludeMap[symbol] {
			filtered = append(filtered, symbol)
		}
	}

	return filtered
}

// convertToEntity converts Binance market data to domain entity
func (c *Collector) convertToEntity(data *binance.MarketData) *entity.MarketData {
	return &entity.MarketData{
		Symbol:                 data.Symbol,
		Timestamp:              data.Timestamp,
		LongAccountRatio:       decimal.NewFromFloat(data.LongAccountRatio),
		ShortAccountRatio:      decimal.NewFromFloat(data.ShortAccountRatio),
		LongPositionRatio:      decimal.NewFromFloat(data.LongPositionRatio),
		ShortPositionRatio:     decimal.NewFromFloat(data.ShortPositionRatio),
		PositionRatioAvailable: data.PositionRatioAvailable,
		DataQualityScore:       data.DataQualityScore,
		Price:                  decimal.NewFromFloat(data.Price),
		Volume24h:              decimal.NewFromFloat(data.Volume24h),
		OpenInterest:           decimal.NewFromFloat(data.OpenInterest),
		FundingRate:            decimal.NewFromFloat(data.FundingRate),
	}
}

// GetCollectionStatus returns the current collection status
func (c *Collector) GetCollectionStatus(ctx context.Context) (map[string]interface{}, error) {
	repo := *c.marketDataRepo

	count, err := repo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get count: %w", err)
	}

	// Get latest data point timestamp
	latestData, err := repo.GetLatestForAllSymbols(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest data: %w", err)
	}

	var latestTimestamp time.Time
	if len(latestData) > 0 {
		latestTimestamp = latestData[0].Timestamp
	}

	status := map[string]interface{}{
		"enabled":           c.config.Enabled,
		"total_data_points": count,
		"symbols_tracked":   len(latestData),
		"latest_collection": latestTimestamp,
	}

	return status, nil
}
