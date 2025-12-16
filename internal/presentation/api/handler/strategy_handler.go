package handler

import (
	"ContractAnalysis/internal/domain/service" // Add this import
	"ContractAnalysis/internal/presentation/api/dto"
	"ContractAnalysis/pkg/utils"

	"github.com/gin-gonic/gin"
)

// StrategyHandler handles strategy related requests
type StrategyHandler struct {
	strategies []service.Strategy // Change from config.StrategiesConfig
}

// NewStrategyHandler creates a new strategy handler
func NewStrategyHandler(strategies []service.Strategy) *StrategyHandler { // Change parameter type
	return &StrategyHandler{
		strategies: strategies, // Assign the slice
	}
}

// GetStrategies returns the list of available strategies
func (h *StrategyHandler) GetStrategies(c *gin.Context) {
	var strategyResponses []dto.StrategyResponse // Renamed to avoid confusion with h.strategies

	for _, s := range h.strategies {
		strategyResponses = append(strategyResponses, dto.StrategyResponse{
			Key:         s.Key(),  // Use s.Key()
			Name:        s.Name(), // Use s.Name()
			Enabled:     s.IsEnabled(),
			Description: s.Name(), // Assuming description is the name for now, or might need another field in Strategy interface
		})
	}

	utils.SuccessResponse(c, 200, "Strategies fetched successfully", strategyResponses)
}
