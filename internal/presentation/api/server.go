package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"ContractAnalysis/internal/domain/repository"
	"ContractAnalysis/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server represents the HTTP API server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	config     ServerConfig
	deps       Dependencies
	logger     *logger.Logger
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Dependencies holds all dependencies for the API server
type Dependencies struct {
	SignalRepo     repository.SignalRepository
	StatisticsRepo repository.StatisticsRepository
	MarketDataRepo repository.MarketDataRepository
	PairRepo       repository.TradingPairRepository
}

// NewServer creates a new API server
func NewServer(config ServerConfig, deps Dependencies, log *logger.Logger, version string) *Server {
	router := SetupRouter(deps, log, version)

	server := &Server{
		router: router,
		config: config,
		deps:   deps,
		logger: log,
	}

	return server
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	s.logger.Info("API server starting", zap.String("address", addr))

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down API server")

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
	}

	s.logger.Info("API server stopped")
	return nil
}
