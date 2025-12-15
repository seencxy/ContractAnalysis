package notification

import (
	"context"

	"ContractAnalysis/internal/domain/entity"
)

// EventType represents the type of notification event
type EventType string

const (
	EventSignalGenerated  EventType = "signal_generated"
	EventSignalConfirmed  EventType = "signal_confirmed"
	EventSignalInvalidated EventType = "signal_invalidated"
	EventSignalOutcome    EventType = "signal_outcome"
	EventSystemError      EventType = "system_error"
)

// Notification represents a notification message
type Notification struct {
	EventType EventType
	Signal    *entity.Signal
	Outcome   *entity.SignalOutcome
	Message   string
	Metadata  map[string]interface{}
}

// Notifier defines the interface for all notifiers
type Notifier interface {
	// Name returns the notifier name
	Name() string

	// IsEnabled returns whether the notifier is enabled
	IsEnabled() bool

	// ShouldNotify checks if this notifier should handle the event
	ShouldNotify(eventType EventType) bool

	// Notify sends a notification
	Notify(ctx context.Context, notification *Notification) error
}

// NotificationDispatcher manages multiple notifiers
type NotificationDispatcher struct {
	notifiers []Notifier
}

// NewNotificationDispatcher creates a new notification dispatcher
func NewNotificationDispatcher(notifiers []Notifier) *NotificationDispatcher {
	return &NotificationDispatcher{
		notifiers: notifiers,
	}
}

// Notify sends a notification to all enabled notifiers
func (d *NotificationDispatcher) Notify(ctx context.Context, notification *Notification) error {
	for _, notifier := range d.notifiers {
		if !notifier.IsEnabled() {
			continue
		}

		if !notifier.ShouldNotify(notification.EventType) {
			continue
		}

		if err := notifier.Notify(ctx, notification); err != nil {
			// Log error but continue with other notifiers
			continue
		}
	}

	return nil
}

// NotifySignalGenerated sends a notification when a signal is generated
func (d *NotificationDispatcher) NotifySignalGenerated(ctx context.Context, signal *entity.Signal) error {
	return d.Notify(ctx, &Notification{
		EventType: EventSignalGenerated,
		Signal:    signal,
		Message:   "New trading signal generated",
	})
}

// NotifySignalConfirmed sends a notification when a signal is confirmed
func (d *NotificationDispatcher) NotifySignalConfirmed(ctx context.Context, signal *entity.Signal) error {
	return d.Notify(ctx, &Notification{
		EventType: EventSignalConfirmed,
		Signal:    signal,
		Message:   "Trading signal confirmed",
	})
}

// NotifySignalInvalidated sends a notification when a signal is invalidated
func (d *NotificationDispatcher) NotifySignalInvalidated(ctx context.Context, signal *entity.Signal) error {
	return d.Notify(ctx, &Notification{
		EventType: EventSignalInvalidated,
		Signal:    signal,
		Message:   "Trading signal invalidated",
	})
}

// NotifySignalOutcome sends a notification when a signal outcome is calculated
func (d *NotificationDispatcher) NotifySignalOutcome(ctx context.Context, signal *entity.Signal, outcome *entity.SignalOutcome) error {
	return d.Notify(ctx, &Notification{
		EventType: EventSignalOutcome,
		Signal:    signal,
		Outcome:   outcome,
		Message:   "Trading signal outcome",
	})
}

// NotifySystemError sends a notification when a system error occurs
func (d *NotificationDispatcher) NotifySystemError(ctx context.Context, message string, metadata map[string]interface{}) error {
	return d.Notify(ctx, &Notification{
		EventType: EventSystemError,
		Message:   message,
		Metadata:  metadata,
	})
}
