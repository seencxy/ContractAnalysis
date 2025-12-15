package redis

import (
	"context"
	"fmt"

	"ContractAnalysis/config"
	"ContractAnalysis/internal/infrastructure/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// NewConnection creates a new Redis client connection
func NewConnection(cfg config.RedisConfig) (*redis.Client, error) {
	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	logger.Info("Successfully connected to Redis",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Int("db", cfg.DB),
	)

	return client, nil
}

// Close closes the Redis client connection
func Close(client *redis.Client) error {
	return client.Close()
}
