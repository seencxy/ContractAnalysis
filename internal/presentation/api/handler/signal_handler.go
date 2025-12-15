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

	// Get signals based on status filter
	var signals []*entity.Signal
	var err error

	if req.Status != "" {
		status := entity.SignalStatus(req.Status)
		signals, err = h.signalRepo.GetByStatus(ctx, status, 1000) // Get more than needed
	} else {
		signals, err = h.signalRepo.GetAll(ctx)
	}

	if err != nil {
		h.logger.Error("Failed to get signals", zap.Error(err))
		apiErr := apierrors.NewDatabaseError("Failed to retrieve signals")
		utils.ErrorResponse(c, apiErr)
		return
	}

	// Filter by symbol if provided
	if req.Symbol != "" {
		filtered := make([]*entity.Signal, 0)
		for _, signal := range signals {
			if signal.Symbol == req.Symbol {
				filtered = append(filtered, signal)
			}
		}
		signals = filtered
	}

	// Filter by strategy if provided
	if req.Strategy != "" {
		filtered := make([]*entity.Signal, 0)
		for _, signal := range signals {
			if signal.StrategyName == req.Strategy {
				filtered = append(filtered, signal)
			}
		}
		signals = filtered
	}

	// Filter by type if provided
	if req.Type != "" {
		signalType := entity.SignalType(req.Type)
		filtered := make([]*entity.Signal, 0)
		for _, signal := range signals {
			if signal.Type == signalType {
				filtered = append(filtered, signal)
			}
		}
		signals = filtered
	}

	// Filter by time range if provided
	if req.StartTime != nil || req.EndTime != nil {
		filtered := make([]*entity.Signal, 0)
		for _, signal := range signals {
			if req.StartTime != nil && signal.GeneratedAt.Before(*req.StartTime) {
				continue
			}
			if req.EndTime != nil && signal.GeneratedAt.After(*req.EndTime) {
				continue
			}
			filtered = append(filtered, signal)
		}
		signals = filtered
	}

	// Calculate total before pagination
	total := len(signals)

	// Apply pagination
	start := pagination.Offset
	end := pagination.Offset + pagination.Limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedSignals := signals[start:end]

	// Get outcomes for CLOSED signals
	closedSignalIDs := make([]string, 0)
	for _, signal := range paginatedSignals {
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
	response := make([]*dto.SignalResponse, 0, len(paginatedSignals))
	for _, signal := range paginatedSignals {
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
