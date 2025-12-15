package handler

import (
	"net/http"
	"time"

	"ContractAnalysis/internal/presentation/api/dto"
	"ContractAnalysis/pkg/utils"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	version string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		version: version,
	}
}

// Check handles GET /api/v1/health
func (h *HealthHandler) Check(c *gin.Context) {
	response := &dto.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   h.version,
	}

	utils.SuccessResponse(c, http.StatusOK, "success", response)
}
