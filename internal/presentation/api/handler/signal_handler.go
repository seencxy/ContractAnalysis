package handler

import (
	"net/http"

	"ContractAnalysis/internal/domain/entity"
	"ContractAnalysis/internal/domain/repository"
	"ContractAnalysis/internal/infrastructure/logger"
	"ContractAnalysis/internal/presentation/api/dto"
	"ContractAnalysis/internal/presentation/api/serializer"
	apierrors "ContractAnalysis/pkg/errors"
	"ContractAnalysis/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignalHandler handles signal-related requests
type SignalHandler struct {
	signalRepo repository.SignalRepository
	logger     *logger.Logger
}

// NewSignalHandler creates a new signal handler
func NewSignalHandler(signalRepo repository.SignalRepository, log *logger.Logger) *SignalHandler {
	return &SignalHandler{
		signalRepo: signalRepo,
		logger:     log,
	}
}

// GetSignals handles GET /api/v1/signals
func (h *SignalHandler) GetSignals(c *gin.Context) {
	var req dto.SignalListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		apiErr := apierrors.NewBadRequestError("Invalid query parameters", err.Error())
		utils.ErrorResponse(c, apiErr)
		return
	}

	// Parse pagination
	pagination, apiErr := utils.ParsePaginationParams(c)
	if apiErr != nil {
		utils.ErrorResponse(c, apiErr)
		return
	}

	ctx := c.Request.Context()

	// Construct filters for the repository
	filters := repository.SignalFilterParams{
		Status:    req.Status,
		Symbol:    req.Symbol,
		Strategy:  req.Strategy,
		Type:      req.Type,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Get signals based on all filters
	signals, total, err := h.signalRepo.GetSignalsWithFilters(ctx, filters, pagination.Offset, pagination.Limit)
	if err != nil {
		h.logger.Error("Failed to get signals with filters", zap.Error(err))
		apiErr := apierrors.NewDatabaseError("Failed to retrieve signals")
		utils.ErrorResponse(c, apiErr)
		return
	}

	// Get outcomes for CLOSED signals
	closedSignalIDs := make([]string, 0)
	for _, signal := range signals {
		if signal.Status == entity.SignalStatusClosed {
			closedSignalIDs = append(closedSignalIDs, signal.SignalID)
		}
	}

	// Fetch outcomes if there are CLOSED signals
	var outcomeMap map[string]*entity.SignalOutcome
	if len(closedSignalIDs) > 0 {
		outcomeMap, err = h.signalRepo.GetOutcomesBySignalIDs(ctx, closedSignalIDs)
		if err != nil {
			h.logger.Warn("Failed to get outcomes for closed signals", zap.Error(err))
			outcomeMap = make(map[string]*entity.SignalOutcome)
		}
	} else {
		outcomeMap = make(map[string]*entity.SignalOutcome)
	}

	// Serialize response with outcomes
	response := make([]*dto.SignalResponse, 0, len(signals))
	for _, signal := range signals {
		outcome := outcomeMap[signal.SignalID]
		response = append(response, serializer.ToSignalResponseWithOutcome(signal, outcome))
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "success", response, pagination.Page, pagination.Limit, total)
}

// GetSignalByID handles GET /api/v1/signals/:id
func (h *SignalHandler) GetSignalByID(c *gin.Context) {
	signalID := c.Param("id")
	ctx := c.Request.Context()

	signal, err := h.signalRepo.GetByID(ctx, signalID)
	if err != nil {
		h.logger.Error("Failed to get signal", zap.String("signal_id", signalID), zap.Error(err))
		apiErr := apierrors.NewNotFoundError("Signal not found")
		utils.ErrorResponse(c, apiErr)
		return
	}

	// Get outcome if signal is CLOSED
	var outcome *entity.SignalOutcome
	if signal.Status == entity.SignalStatusClosed {
		outcome, err = h.signalRepo.GetOutcome(ctx, signalID)
		if err != nil {
			h.logger.Warn("Failed to get outcome for closed signal", zap.String("signal_id", signalID), zap.Error(err))
		}
	}

	response := serializer.ToSignalResponseWithOutcome(signal, outcome)
	utils.SuccessResponse(c, http.StatusOK, "success", response)
}

// GetSignalTracking handles GET /api/v1/signals/:id/tracking
func (h *SignalHandler) GetSignalTracking(c *gin.Context) {
	signalID := c.Param("id")
	ctx := c.Request.Context()

	trackings, err := h.signalRepo.GetAllTracking(ctx, signalID)
	if err != nil {
		h.logger.Error("Failed to get signal tracking", zap.String("signal_id", signalID), zap.Error(err))
		apiErr := apierrors.NewDatabaseError("Failed to retrieve tracking data")
		utils.ErrorResponse(c, apiErr)
		return
	}

	response := serializer.ToSignalTrackingListResponse(trackings)
	utils.SuccessResponse(c, http.StatusOK, "success", response)
}

// GetSignalKlines handles GET /api/v1/signals/:id/klines
func (h *SignalHandler) GetSignalKlines(c *gin.Context) {
	signalID := c.Param("id")
	ctx := c.Request.Context()

	klines, err := h.signalRepo.GetKlineTrackingBySignal(ctx, signalID)
	if err != nil {
		h.logger.Error("Failed to get signal klines", zap.String("signal_id", signalID), zap.Error(err))
		apiErr := apierrors.NewDatabaseError("Failed to retrieve kline data")
		utils.ErrorResponse(c, apiErr)
		return
	}

	response := serializer.ToSignalKlineTrackingListResponse(klines)
	utils.SuccessResponse(c, http.StatusOK, "success", response)
}

// GetActiveSignals handles GET /api/v1/signals/active
func (h *SignalHandler) GetActiveSignals(c *gin.Context) {
	ctx := c.Request.Context()

	signals, err := h.signalRepo.GetActiveSignals(ctx)
	if err != nil {
		h.logger.Error("Failed to get active signals", zap.Error(err))
		apiErr := apierrors.NewDatabaseError("Failed to retrieve active signals")
		utils.ErrorResponse(c, apiErr)
		return
	}

	response := serializer.ToSignalListResponse(signals)
	utils.SuccessResponse(c, http.StatusOK, "success", response)
}
