package handler

import (
	"context"
	"net/http"
	"time"

	"ContractAnalysis/internal/domain/entity"
	"ContractAnalysis/internal/domain/repository"
	"ContractAnalysis/internal/infrastructure/logger"
	"ContractAnalysis/internal/presentation/api/dto"
	"ContractAnalysis/internal/presentation/api/serializer"
	apierrors "ContractAnalysis/pkg/errors"
	"ContractAnalysis/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// StatisticsHandler handles statistics-related requests
type StatisticsHandler struct {
	statisticsRepo repository.StatisticsRepository
	signalRepo     repository.SignalRepository
	logger         *logger.Logger
}

// NewStatisticsHandler creates a new statistics handler
func NewStatisticsHandler(statsRepo repository.StatisticsRepository, signalRepo repository.SignalRepository, log *logger.Logger) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsRepo: statsRepo,
		signalRepo:     signalRepo,
		logger:         log,
	}
}

// GetOverview handles GET /api/v1/statistics/overview
func (h *StatisticsHandler) GetOverview(c *gin.Context) {
	ctx := c.Request.Context()

	// Calculate overview statistics
	overview, err := h.calculateOverviewStatistics(ctx)
	if err != nil {
		h.logger.Error("Failed to calculate overview statistics", zap.Error(err))
		apiErr := apierrors.NewDatabaseError("Failed to retrieve overview statistics")
		utils.ErrorResponse(c, apiErr)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", overview)
}

// GetStrategies handles GET /api/v1/statistics/strategies
func (h *StatisticsHandler) GetStrategies(c *gin.Context) {
	var req dto.StatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		apiErr := apierrors.NewBadRequestError("Invalid query parameters", err.Error())
		utils.ErrorResponse(c, apiErr)
		return
	}

	ctx := c.Request.Context()

	// Default period to "all" if not specified
	period := req.Period
	if period == "" {
		period = "all"
	}

	// Optional strategy filter
	var strategyFilter *string
	if req.StrategyName != "" {
		strategyFilter = &req.StrategyName
	}

	// Get statistics for the period, with optional strategy filter
	stats, err := h.statisticsRepo.GetByPeriodAndStrategy(ctx, period, strategyFilter)
	if err != nil {
		h.logger.Error("Failed to get strategy statistics", zap.String("period", period), zap.Error(err), zap.Stringp("strategy", strategyFilter))
		apiErr := apierrors.NewDatabaseError("Failed to retrieve statistics")
		utils.ErrorResponse(c, apiErr)
		return
	}

	responses := serializer.ToStatisticsListResponse(stats)

	utils.SuccessResponse(c, http.StatusOK, "success", responses)
}

// GetSymbols handles GET /api/v1/statistics/symbols
func (h *StatisticsHandler) GetSymbols(c *gin.Context) {
	var req dto.StatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		apiErr := apierrors.NewBadRequestError("Invalid query parameters", err.Error())
		utils.ErrorResponse(c, apiErr)
		return
	}

	ctx := c.Request.Context()

	// Default period to "all" if not specified
	period := req.Period
	if period == "" {
		period = "all"
	}

	// Get statistics for the period
	stats, err := h.statisticsRepo.GetByPeriod(ctx, period)
	if err != nil {
		h.logger.Error("Failed to get symbol statistics", zap.String("period", period), zap.Error(err))
		apiErr := apierrors.NewDatabaseError("Failed to retrieve statistics")
		utils.ErrorResponse(c, apiErr)
		return
	}

	// Filter to only include symbol-specific stats
	filtered := make([]*repository.StrategyStatistics, 0)
	for _, stat := range stats {
		if stat.Symbol != nil {
			filtered = append(filtered, stat)
		}
	}

	responses := serializer.ToStatisticsListResponse(filtered)

	utils.SuccessResponse(c, http.StatusOK, "success", responses)
}

// GetHistory handles GET /api/v1/statistics/history
func (h *StatisticsHandler) GetHistory(c *gin.Context) {
	var req dto.StatisticsHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		apiErr := apierrors.NewBadRequestError("Invalid request parameters", err.Error())
		utils.ErrorResponse(c, apiErr)
		return
	}

	// Validate time range
	if req.StartTime == nil || req.EndTime == nil {
		apiErr := apierrors.NewBadRequestError("start_time and end_time are required", "")
		utils.ErrorResponse(c, apiErr)
		return
	}

	if req.EndTime.Before(*req.StartTime) {
		apiErr := apierrors.NewBadRequestError("end_time must be after start_time", "")
		utils.ErrorResponse(c, apiErr)
		return
	}

	ctx := c.Request.Context()

	// Optional filters
	var strategyFilter *string
	if req.StrategyName != "" {
		strategyFilter = &req.StrategyName
	}

	var symbolFilter *string
	if req.Symbol != "" {
		symbolFilter = &req.Symbol
	}

	// Get historical statistics
	stats, err := h.statisticsRepo.GetByTimeRange(
		ctx,
		*req.StartTime,
		*req.EndTime,
		strategyFilter,
		symbolFilter,
	)

	if err != nil {
		h.logger.Error("Failed to get historical statistics",
			zap.Time("start_time", *req.StartTime),
			zap.Time("end_time", *req.EndTime),
			zap.Error(err))
		apiErr := apierrors.NewDatabaseError("Failed to retrieve historical statistics")
		utils.ErrorResponse(c, apiErr)
		return
	}

	responses := serializer.ToStatisticsListResponse(stats)

	utils.SuccessResponse(c, http.StatusOK, "success", responses)
}

// calculateOverviewStatistics calculates overview statistics for dashboard
func (h *StatisticsHandler) calculateOverviewStatistics(ctx context.Context) (*dto.OverviewStatisticsResponse, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Get today's signals
	todaySignals, err := h.signalRepo.GetSignalsInTimeRange(ctx, todayStart, now)
	if err != nil {
		return nil, err
	}

	// Get active signals (PENDING, CONFIRMED, TRACKING)
	activeSignals, err := h.signalRepo.GetActiveSignals(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate status distribution from all signals today
	statusDistribution := &dto.SignalStatusDistribution{}
	for _, signal := range todaySignals {
		switch signal.Status {
		case entity.SignalStatusPending:
			statusDistribution.Pending++
		case entity.SignalStatusConfirmed:
			statusDistribution.Confirmed++
		case entity.SignalStatusTracking:
			statusDistribution.Tracking++
		case entity.SignalStatusClosed:
			statusDistribution.Closed++
		case entity.SignalStatusInvalidated:
			statusDistribution.Invalidated++
		}
	}

	// Get 24h statistics from the statistics table
	stats24h, err := h.statisticsRepo.GetByPeriod(ctx, "24h")
	if err != nil {
		h.logger.Error("Failed to get 24h statistics", zap.Error(err))
		return nil, err
	}

	h.logger.Info("Retrieved 24h statistics", zap.Int("count", len(stats24h)))

	// If no 24h data, try to get "all" period as fallback
	if len(stats24h) == 0 {
		h.logger.Warn("No 24h statistics found, trying 'all' period as fallback")
		stats24h, err = h.statisticsRepo.GetByPeriod(ctx, "all")
		if err != nil {
			h.logger.Error("Failed to get 'all' statistics", zap.Error(err))
			return nil, err
		}
		h.logger.Info("Retrieved 'all' statistics as fallback", zap.Int("count", len(stats24h)))
	}

	// Initialize response with defaults
	zeroStr := "0"
	response := &dto.OverviewStatisticsResponse{
		TotalSignalsToday:   len(todaySignals),
		ActiveSignals:       len(activeSignals),
		OverallWinRate24h:   &zeroStr,
		AvgReturnPct24h:     &zeroStr,
		TopPerformingPair:   "-",
		WorstPerformingPair: "-",
		StatusDistribution:  statusDistribution,
	}

	// Calculate overall 24h metrics from statistics
	if len(stats24h) > 0 {
		var totalSignals int
		var totalProfitable int
		var totalReturn decimal.Decimal
		pairReturns := make(map[string]decimal.Decimal)
		pairCounts := make(map[string]int)

		h.logger.Info("Processing statistics", zap.Int("stat_count", len(stats24h)))

		// Aggregate all statistics (whether overall or per-symbol)
		var recordsWithSignals int
		var recordsWithoutAvgPct int

		for _, stat := range stats24h {
			signalCount := stat.ProfitableSignals + stat.LosingSignals

			if signalCount > 0 {
				recordsWithSignals++

				// Accumulate for overall metrics
				totalSignals += signalCount
				totalProfitable += stat.ProfitableSignals

				// Calculate weighted return
				if stat.AvgProfitPct != nil && stat.AvgLossPct != nil {
					profitWeight := decimal.NewFromInt(int64(stat.ProfitableSignals))
					lossWeight := decimal.NewFromInt(int64(stat.LosingSignals))

					profitContribution := stat.AvgProfitPct.Mul(profitWeight)
					lossContribution := stat.AvgLossPct.Mul(lossWeight).Neg() // Loss is negative

					totalReturn = totalReturn.Add(profitContribution).Add(lossContribution)
				} else {
					recordsWithoutAvgPct++
				}

				// Track per-symbol performance
				if stat.Symbol != nil {
					symbol := *stat.Symbol

					if stat.AvgProfitPct != nil && stat.AvgLossPct != nil {
						profitWeight := decimal.NewFromInt(int64(stat.ProfitableSignals))
						lossWeight := decimal.NewFromInt(int64(stat.LosingSignals))

						profitContribution := stat.AvgProfitPct.Mul(profitWeight)
						lossContribution := stat.AvgLossPct.Mul(lossWeight).Neg()

						symbolReturn := profitContribution.Add(lossContribution)
						pairReturns[symbol] = pairReturns[symbol].Add(symbolReturn)
						pairCounts[symbol] += signalCount
					}
				}
			}
		}

		h.logger.Info("Aggregated statistics",
			zap.Int("total_records", len(stats24h)),
			zap.Int("records_with_signals", recordsWithSignals),
			zap.Int("records_without_avg_pct", recordsWithoutAvgPct),
			zap.Int("total_signals", totalSignals),
			zap.Int("total_profitable", totalProfitable))

		// Calculate overall win rate
		if totalSignals > 0 {
			winRate := decimal.NewFromInt(int64(totalProfitable)).Div(decimal.NewFromInt(int64(totalSignals))).Mul(decimal.NewFromInt(100))
			winRateStr := winRate.String()
			response.OverallWinRate24h = &winRateStr
			h.logger.Info("Calculated win rate",
				zap.Int("total_signals", totalSignals),
				zap.Int("profitable", totalProfitable),
				zap.String("win_rate", winRateStr))
		} else {
			h.logger.Warn("No signals to calculate win rate")
		}

		// Calculate average return
		if totalSignals > 0 {
			avgReturn := totalReturn.Div(decimal.NewFromInt(int64(totalSignals)))
			avgReturnStr := avgReturn.String()
			response.AvgReturnPct24h = &avgReturnStr
			h.logger.Info("Calculated average return",
				zap.String("avg_return", avgReturnStr))
		}

		// Find top and worst performing pairs
		if len(pairReturns) > 0 {
			var topPair, worstPair string
			var topReturn, worstReturn decimal.Decimal
			first := true

			for symbol, totalRet := range pairReturns {
				avgRet := totalRet.Div(decimal.NewFromInt(int64(pairCounts[symbol])))

				if first {
					topPair = symbol
					worstPair = symbol
					topReturn = avgRet
					worstReturn = avgRet
					first = false
				} else {
					if avgRet.GreaterThan(topReturn) {
						topPair = symbol
						topReturn = avgRet
					}
					if avgRet.LessThan(worstReturn) {
						worstPair = symbol
						worstReturn = avgRet
					}
				}
			}

			response.TopPerformingPair = topPair
			response.WorstPerformingPair = worstPair
			h.logger.Info("Calculated top/worst pairs",
				zap.String("top", topPair),
				zap.String("worst", worstPair))
		}
	} else {
		h.logger.Warn("No statistics data available for overview calculation")
	}

	h.logger.Info("Overview statistics calculated",
		zap.Int("today_signals", response.TotalSignalsToday),
		zap.Int("active_signals", response.ActiveSignals),
		zap.Bool("has_win_rate", response.OverallWinRate24h != nil),
		zap.Bool("has_avg_return", response.AvgReturnPct24h != nil))

	return response, nil
}
