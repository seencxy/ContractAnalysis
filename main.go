package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ContractAnalysis/config"
	"ContractAnalysis/internal/domain/repository"
	"ContractAnalysis/internal/domain/service"
	"ContractAnalysis/internal/infrastructure/binance"
	"ContractAnalysis/internal/infrastructure/logger"
	"ContractAnalysis/internal/infrastructure/notification"
	mysqlRepo "ContractAnalysis/internal/infrastructure/persistence/mysql"
	redisConn "ContractAnalysis/internal/infrastructure/persistence/redis"
	"ContractAnalysis/internal/infrastructure/scheduler"
	"ContractAnalysis/internal/presentation/api"
	"ContractAnalysis/internal/usecase"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(logger.Config{
		Level:      cfg.Logging.Level,
		Format:     cfg.Logging.Format,
		Output:     cfg.Logging.Output,
		FilePath:   cfg.Logging.File.Path,
		MaxSize:    cfg.Logging.File.MaxSize,
		MaxBackups: cfg.Logging.File.MaxBackups,
		MaxAge:     cfg.Logging.File.MaxAge,
		Compress:   cfg.Logging.File.Compress,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Set global logger
	logger.SetGlobal(log)

	log.Info("Starting Binance Futures Analysis System",
		zap.String("version", cfg.App.Version),
		zap.String("environment", cfg.App.Environment),
	)

	// Initialize database connections
	log.Info("Connecting to MySQL...")
	db, err := mysqlRepo.NewConnection(cfg.Database.MySQL)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to MySQL")
	}

	log.Info("Connecting to Redis...")
	redisClient, err := redisConn.NewConnection(cfg.Database.Redis)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to Redis")
	}
	defer redisClient.Close()

	// Initialize Binance client
	log.Info("Initializing Binance client...")
	binanceClient, err := binance.NewClient(cfg.Binance)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize Binance client")
	}

	// Initialize repositories
	tradingPairRepo := mysqlRepo.NewTradingPairRepository(db)
	marketDataRepoImpl := mysqlRepo.NewMarketDataRepository(db)
	marketDataRepo := repository.MarketDataRepository(marketDataRepoImpl)
	signalRepoImpl := mysqlRepo.NewSignalRepository(db)
	signalRepo := repository.SignalRepository(signalRepoImpl)
	statisticsRepo := mysqlRepo.NewStatisticsRepository(db)

	// Initialize strategies
	var strategies []service.Strategy

	if cfg.Strategies.Minority.Enabled {
		minorityStrategy := service.NewMinorityStrategy(service.MinorityStrategyConfig{
			BaseConfig: service.StrategyConfig{
				Name:              cfg.Strategies.Minority.Name,
				Enabled:           cfg.Strategies.Minority.Enabled,
				ConfirmationHours: cfg.Strategies.Minority.ConfirmationHours,
				TrackingHours:     cfg.Strategies.Minority.TrackingHours,
				ProfitTargetPct:   cfg.Strategies.Minority.ProfitTargetPct,
				StopLossPct:       cfg.Strategies.Minority.StopLossPct,
			},
			MinRatioDifference:              cfg.Strategies.Minority.MinRatioDifference,
			GenerateLongWhenShortRatioAbove: cfg.Strategies.Minority.GenerateLongWhenShortRatioAbove,
			GenerateShortWhenLongRatioAbove: cfg.Strategies.Minority.GenerateShortWhenLongRatioAbove,
		})
		strategies = append(strategies, minorityStrategy)
		log.Info("Minority strategy enabled")
	}

	if cfg.Strategies.Whale.Enabled {
		whaleStrategy := service.NewWhaleStrategy(service.WhaleStrategyConfig{
			BaseConfig: service.StrategyConfig{
				Name:              cfg.Strategies.Whale.Name,
				Enabled:           cfg.Strategies.Whale.Enabled,
				ConfirmationHours: cfg.Strategies.Whale.ConfirmationHours,
				TrackingHours:     cfg.Strategies.Whale.TrackingHours,
				ProfitTargetPct:   cfg.Strategies.Whale.ProfitTargetPct,
				StopLossPct:       cfg.Strategies.Whale.StopLossPct,
			},
			MinRatioDifference:     cfg.Strategies.Whale.MinRatioDifference,
			WhalePositionThreshold: cfg.Strategies.Whale.WhalePositionThreshold,
			MinDivergence:          cfg.Strategies.Whale.MinDivergence,
		})
		strategies = append(strategies, whaleStrategy)
		log.Info("Whale strategy enabled")
	}

	if cfg.Strategies.SmartMoney.Enabled {
		smartMoneyStrategy := service.NewSmartMoneyStrategy(service.SmartMoneyStrategyConfig{
			BaseConfig: service.StrategyConfig{
				Name:              cfg.Strategies.SmartMoney.Name,
				Enabled:           cfg.Strategies.SmartMoney.Enabled,
				ConfirmationHours: cfg.Strategies.SmartMoney.ConfirmationHours,
				TrackingHours:     cfg.Strategies.SmartMoney.TrackingHours,
				ProfitTargetPct:   cfg.Strategies.SmartMoney.ProfitTargetPct,
				StopLossPct:       cfg.Strategies.SmartMoney.StopLossPct,
			},
			MinLongAccountRatio: cfg.Strategies.SmartMoney.MinLongAccountRatio,
			LookbackPeriod:      cfg.Strategies.SmartMoney.LookbackPeriod,
			KlineInterval:       cfg.Strategies.SmartMoney.KlineInterval,
		}, binanceClient) // Use binanceClient as klineRepo
		strategies = append(strategies, smartMoneyStrategy)
		log.Info("Smart Money strategy enabled")
	}

	log.Info("Strategies initialized", zap.Int("count", len(strategies)))

	// Initialize notification system
	var notifiers []notification.Notifier

	if cfg.Notifications.Console.Enabled {
		consoleNotifier := notification.NewConsoleNotifier(cfg.Notifications.Console)
		notifiers = append(notifiers, consoleNotifier)
		log.Info("Console notifier enabled")
	}

	notificationDispatcher := notification.NewNotificationDispatcher(notifiers)

	// Initialize use cases
	collector := usecase.NewCollector(
		binanceClient,
		&marketDataRepo,
		tradingPairRepo,
		cfg.Collection,
	)

	analyzer := usecase.NewAnalyzer(
		strategies,
		&marketDataRepo,
		&signalRepo,
		tradingPairRepo,
		cfg.Strategies.Global,
	)

	tracker := usecase.NewTracker(
		binanceClient,
		&signalRepo,
	)

	statisticsCalculator := usecase.NewStatisticsCalculator(
		&signalRepo,
		statisticsRepo,
		cfg.Statistics,
	)

	// Initialize statistics monitor
	statisticsMonitor := usecase.NewStatisticsMonitor(
		statisticsRepo,
		cfg.Statistics.Monitoring,
	)

	// Initialize API server
	apiServer := api.NewServer(
		api.ServerConfig{
			Host:         cfg.Server.Host,
			Port:         cfg.Server.Port,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
		api.Dependencies{
			SignalRepo:       signalRepo,
			StatsRepo:        statisticsRepo,
			MarketDataRepo:   marketDataRepo,
			TradingPairRepo:  tradingPairRepo,
			StrategiesConfig: cfg.Strategies,
			Strategies:       strategies, // Add this line
		},
		log,
		cfg.App.Version,
	)

	// Start API server in goroutine
	go func() {
		log.Info("Starting API server", zap.Int("port", cfg.Server.Port))
		if err := apiServer.Start(); err != nil {
			log.WithError(err).Fatal("Failed to start API server")
		}
	}()

	// Initialize scheduler
	sched := scheduler.NewScheduler(
		collector,
		analyzer,
		tracker,
		statisticsCalculator,
		statisticsMonitor,
		notificationDispatcher,
	)

	// Add scheduled jobs
	if cfg.Collection.Enabled {
		// Data collection job (every hour by default)
		if err = sched.AddCollectionJob(cfg.Collection.Interval); err != nil {
			log.WithError(err).Fatal("Failed to add collection job")
		}

		// Signal analysis job (every hour at minute 5)
		if err = sched.AddAnalysisJob("0 5 * * * *"); err != nil {
			log.WithError(err).Fatal("Failed to add analysis job")
		}
	}

	// Signal tracking job (every 15 minutes)
	if err = sched.AddTrackingJob("0 */15 * * * *"); err != nil {
		log.WithError(err).Fatal("Failed to add tracking job")
	}

	// Kline tracking job (every hour at minute 5)
	if err = sched.AddKlineTrackingJob("0 5 * * * *"); err != nil {
		log.WithError(err).Fatal("Failed to add kline tracking job")
	}

	// Statistics calculation job (every 6 hours)
	if err = sched.AddStatisticsJob(cfg.Statistics.CalculationInterval); err != nil {
		log.WithError(err).Fatal("Failed to add statistics job")
	}

	// Start scheduler
	sched.Start()

	log.Info("System started successfully")

	// Run initial data collection on startup
	if cfg.Collection.Enabled {
		log.Info("Running initial data collection...")
		ctx := context.Background()
		if err := collector.CollectAll(ctx); err != nil {
			log.WithError(err).Warn("Initial data collection failed, will retry on next scheduled run")
		} else {
			log.Info("Initial data collection completed successfully")

			// Run initial signal analysis after data collection
			log.Info("Running initial signal analysis...")
			signals, err := analyzer.AnalyzeAll(ctx)
			if err != nil {
				log.WithError(err).Warn("Initial signal analysis failed")
			} else {
				log.Info("Initial signal analysis completed", zap.Int("signals_generated", len(signals)))

				// Send notifications for generated signals
				for _, signal := range signals {
					if err := notificationDispatcher.NotifySignalGenerated(ctx, signal); err != nil {
						log.WithError(err).Warn("Failed to send signal notification")
					}
				}
			}
		}
	}

	log.Info("Press Ctrl+C to stop")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down...")

	// Stop scheduler
	sched.Stop()

	// Shutdown API server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Error shutting down API server")
	}

	// Close database connection
	if err := mysqlRepo.Close(db); err != nil {
		log.WithError(err).Error("Error closing database")
	}

	log.Info("Shutdown complete")
}

func init() {
	// Set decimal precision for financial calculations
	decimal.DivisionPrecision = 10
}
