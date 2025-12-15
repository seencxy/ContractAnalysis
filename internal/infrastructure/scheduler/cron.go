package scheduler

import (
	"context"
	"fmt"

	"ContractAnalysis/internal/infrastructure/logger"
	"ContractAnalysis/internal/infrastructure/notification"
	"ContractAnalysis/internal/usecase"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Scheduler manages scheduled jobs
type Scheduler struct {
	cron                 *cron.Cron
	collector            *usecase.Collector
	analyzer             *usecase.Analyzer
	tracker              *usecase.Tracker
	statisticsCalculator *usecase.StatisticsCalculator
	statisticsMonitor    *usecase.StatisticsMonitor
	notifier             *notification.NotificationDispatcher
	logger               *logger.Logger
	ctx                  context.Context
	cancelFunc           context.CancelFunc
}

// NewScheduler creates a new scheduler
func NewScheduler(
	collector *usecase.Collector,
	analyzer *usecase.Analyzer,
	tracker *usecase.Tracker,
	statisticsCalculator *usecase.StatisticsCalculator,
	statisticsMonitor *usecase.StatisticsMonitor,
	notifier *notification.NotificationDispatcher,
) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		cron:                 cron.New(cron.WithSeconds()),
		collector:            collector,
		analyzer:             analyzer,
		tracker:              tracker,
		statisticsCalculator: statisticsCalculator,
		statisticsMonitor:    statisticsMonitor,
		notifier:             notifier,
		logger:               logger.WithComponent("scheduler"),
		ctx:                  ctx,
		cancelFunc:           cancel,
	}
}

// AddCollectionJob adds the data collection job
func (s *Scheduler) AddCollectionJob(schedule string) error {
	_, err := s.cron.AddFunc(schedule, func() {
		s.logger.Info("Running data collection job")

		if err := s.collector.CollectAll(s.ctx); err != nil {
			s.logger.WithError(err).Error("Data collection job failed")
			_ = s.notifier.NotifySystemError(s.ctx, "Data collection failed: "+err.Error(), nil)
			return
		}

		s.logger.Info("Data collection job completed")
	})

	if err != nil {
		return fmt.Errorf("failed to add collection job: %w", err)
	}

	s.logger.Info("Added data collection job", zap.String("schedule", schedule))
	return nil
}

// AddAnalysisJob adds the signal analysis job
func (s *Scheduler) AddAnalysisJob(schedule string) error {
	_, err := s.cron.AddFunc(schedule, func() {
		s.logger.Info("Running signal analysis job")

		// Analyze all symbols
		signals, err := s.analyzer.AnalyzeAll(s.ctx)
		if err != nil {
			s.logger.WithError(err).Error("Signal analysis job failed")
			_ = s.notifier.NotifySystemError(s.ctx, "Signal analysis failed: "+err.Error(), nil)
			return
		}

		// Send notifications for new signals
		for _, signal := range signals {
			if err := s.notifier.NotifySignalGenerated(s.ctx, signal); err != nil {
				s.logger.WithError(err).WithSignalID(signal.SignalID).Warn("Failed to send signal notification")
			}
		}

		// Validate pending signals
		if err := s.analyzer.ValidatePendingSignals(s.ctx); err != nil {
			s.logger.WithError(err).Error("Signal validation failed")
			return
		}

		s.logger.Info("Signal analysis job completed", zap.Int("signals", len(signals)))
	})

	if err != nil {
		return fmt.Errorf("failed to add analysis job: %w", err)
	}

	s.logger.Info("Added signal analysis job", zap.String("schedule", schedule))
	return nil
}

// AddTrackingJob adds the signal tracking job
func (s *Scheduler) AddTrackingJob(schedule string) error {
	_, err := s.cron.AddFunc(schedule, func() {
		s.logger.Info("Running signal tracking job")

		if err := s.tracker.TrackAll(s.ctx); err != nil {
			s.logger.WithError(err).Error("Signal tracking job failed")
			_ = s.notifier.NotifySystemError(s.ctx, "Signal tracking failed: "+err.Error(), nil)
			return
		}

		s.logger.Info("Signal tracking job completed")
	})

	if err != nil {
		return fmt.Errorf("failed to add tracking job: %w", err)
	}

	s.logger.Info("Added signal tracking job", zap.String("schedule", schedule))
	return nil
}

// AddStatisticsJob adds the statistics calculation job
func (s *Scheduler) AddStatisticsJob(schedule string) error {
	_, err := s.cron.AddFunc(schedule, func() {
		s.logger.Info("Running statistics calculation job")

		if err := s.statisticsCalculator.CalculateAll(s.ctx); err != nil {
			s.logger.WithError(err).Error("Statistics calculation job failed")
			_ = s.notifier.NotifySystemError(s.ctx, "Statistics calculation failed: "+err.Error(), nil)
			return
		}

		// Monitor statistics changes if enabled
		if s.statisticsMonitor != nil {
			if err := s.statisticsMonitor.MonitorAllStatistics(s.ctx); err != nil {
				s.logger.WithError(err).Warn("Statistics monitoring failed")
				// Don't fail the job if monitoring fails
			}
		}

		s.logger.Info("Statistics calculation job completed")
	})

	if err != nil {
		return fmt.Errorf("failed to add statistics job: %w", err)
	}

	s.logger.Info("Added statistics calculation job", zap.String("schedule", schedule))
	return nil
}

// AddKlineTrackingJob adds the kline tracking job
func (s *Scheduler) AddKlineTrackingJob(schedule string) error {
	_, err := s.cron.AddFunc(schedule, func() {
		s.logger.Info("Running kline tracking job")

		if err := s.tracker.TrackAllKlines(s.ctx); err != nil {
			s.logger.WithError(err).Error("Kline tracking job failed")
			_ = s.notifier.NotifySystemError(s.ctx, "Kline tracking failed: "+err.Error(), nil)
			return
		}

		s.logger.Info("Kline tracking job completed")
	})

	if err != nil {
		return fmt.Errorf("failed to add kline tracking job: %w", err)
	}

	s.logger.Info("Added kline tracking job", zap.String("schedule", schedule))
	return nil
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.logger.Info("Starting scheduler")
	s.cron.Start()
	s.logger.Info("Scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.logger.Info("Stopping scheduler")
	s.cancelFunc()
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.logger.Info("Scheduler stopped")
}

// GetEntries returns the scheduled job entries
func (s *Scheduler) GetEntries() []cron.Entry {
	return s.cron.Entries()
}
