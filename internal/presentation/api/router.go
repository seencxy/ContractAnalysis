package api

import (
	"ContractAnalysis/internal/infrastructure/logger"
	"ContractAnalysis/internal/presentation/api/handler"
	"ContractAnalysis/internal/presentation/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the HTTP router
func SetupRouter(deps Dependencies, log *logger.Logger, version string) *gin.Engine {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Global middleware
	router.Use(middleware.Recovery(log))
	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS())

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(version)
	signalHandler := handler.NewSignalHandler(deps.SignalRepo, log)
	statisticsHandler := handler.NewStatisticsHandler(deps.StatsRepo, deps.SignalRepo, log)
	strategyHandler := handler.NewStrategyHandler(deps.Strategies)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", healthHandler.Check)

		// Strategies meta
		v1.GET("/strategies", strategyHandler.GetStrategies)

		// Signal routes
		signals := v1.Group("/signals")
		{
			signals.GET("", signalHandler.GetSignals)
			signals.GET("/active", signalHandler.GetActiveSignals)
			signals.GET("/:id", signalHandler.GetSignalByID)
			signals.GET("/:id/tracking", signalHandler.GetSignalTracking)
			signals.GET("/:id/klines", signalHandler.GetSignalKlines)
		}

		// Statistics routes
		statistics := v1.Group("/statistics")
		{
			statistics.GET("/overview", statisticsHandler.GetOverview)
			statistics.GET("/strategies", statisticsHandler.GetStrategies)
			statistics.GET("/symbols", statisticsHandler.GetSymbols)
			statistics.GET("/history", statisticsHandler.GetHistory)
		}
	}

	return router
}
