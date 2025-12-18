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
		apiErr := apierrors.NewValidationError("Invalid query parameters", err.Error())
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
		Status:       req.Status,
		Symbol:       req.Symbol,
		StrategyName: req.StrategyName,
		Type:         req.Type,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
	}

	// Get signals with outcomes using single LEFT JOIN query (optimized)
	signalsWithOutcomes, total, err := h.signalRepo.GetSignalsWithOutcomes(ctx, filters, pagination.Offset, pagination.Limit)
	if err != nil {
		h.logger.Error("Failed to get signals with outcomes", zap.Error(err))
		apiErr := apierrors.NewDatabaseError("Failed to retrieve signals")
		utils.ErrorResponse(c, apiErr)
		return
	}

	// Serialize response with outcomes
	response := make([]*dto.SignalResponse, 0, len(signalsWithOutcomes))
	for _, swo := range signalsWithOutcomes {
		response = append(response, serializer.ToSignalResponseWithOutcome(swo.Signal, swo.Outcome))
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
