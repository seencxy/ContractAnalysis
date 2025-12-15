package notification

import (
	"context"
	"fmt"

	"ContractAnalysis/config"
	"ContractAnalysis/internal/infrastructure/logger"
)

// ConsoleNotifier sends notifications to the console/logs
type ConsoleNotifier struct {
	config config.ConsoleConfig
	logger *logger.Logger
}

// NewConsoleNotifier creates a new console notifier
func NewConsoleNotifier(cfg config.ConsoleConfig) *ConsoleNotifier {
	return &ConsoleNotifier{
		config: cfg,
		logger: logger.WithComponent("console-notifier"),
	}
}

// Name returns the notifier name
func (n *ConsoleNotifier) Name() string {
	return "console"
}

// IsEnabled returns whether the notifier is enabled
func (n *ConsoleNotifier) IsEnabled() bool {
	return n.config.Enabled
}

// ShouldNotify checks if this notifier should handle the event
func (n *ConsoleNotifier) ShouldNotify(eventType EventType) bool {
	for _, event := range n.config.Events {
		if event == string(eventType) {
			return true
		}
	}
	return false
}

// Notify sends a notification to the console
func (n *ConsoleNotifier) Notify(ctx context.Context, notification *Notification) error {
	switch notification.EventType {
	case EventSignalGenerated:
		return n.notifySignalGenerated(notification)
	case EventSignalConfirmed:
		return n.notifySignalConfirmed(notification)
	case EventSignalInvalidated:
		return n.notifySignalInvalidated(notification)
	case EventSignalOutcome:
		return n.notifySignalOutcome(notification)
	case EventSystemError:
		return n.notifySystemError(notification)
	default:
		return fmt.Errorf("unknown event type: %s", notification.EventType)
	}
}

func (n *ConsoleNotifier) notifySignalGenerated(notification *Notification) error {
	signal := notification.Signal
	if signal == nil {
		return fmt.Errorf("signal is nil")
	}

	message := fmt.Sprintf(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🚨 NEW TRADING SIGNAL GENERATED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Signal ID:  %s
Symbol:     %s
Direction:  %s
Strategy:   %s
Price:      %s
Generated:  %s

📊 Market Ratios:
Long/Short (Accounts):  %.2f%% / %.2f%%
Long/Short (Position):  %.2f%% / %.2f%%

📝 Reason:
%s

⏰ Confirmation Period: %d hours
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`,
		signal.SignalID,
		signal.Symbol,
		signal.Type,
		signal.StrategyName,
		signal.PriceAtSignal.String(),
		signal.GeneratedAt.Format("2006-01-02 15:04:05"),
		signal.LongAccountRatio.InexactFloat64(),
		signal.ShortAccountRatio.InexactFloat64(),
		signal.LongPositionRatio.InexactFloat64(),
		signal.ShortPositionRatio.InexactFloat64(),
		signal.Reason,
		int(signal.ConfirmationEnd.Sub(signal.ConfirmationStart).Hours()),
	)

	n.logger.Info(message)
	return nil
}

func (n *ConsoleNotifier) notifySignalConfirmed(notification *Notification) error {
	signal := notification.Signal
	if signal == nil {
		return fmt.Errorf("signal is nil")
	}

	message := fmt.Sprintf(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ TRADING SIGNAL CONFIRMED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Signal ID:  %s
Symbol:     %s
Direction:  %s
Strategy:   %s
Price:      %s
Confirmed:  %s

⚠️  Signal has been confirmed and is now being tracked.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`,
		signal.SignalID,
		signal.Symbol,
		signal.Type,
		signal.StrategyName,
		signal.PriceAtSignal.String(),
		signal.ConfirmedAt.Format("2006-01-02 15:04:05"),
	)

	n.logger.Info(message)
	return nil
}

func (n *ConsoleNotifier) notifySignalInvalidated(notification *Notification) error {
	signal := notification.Signal
	if signal == nil {
		return fmt.Errorf("signal is nil")
	}

	message := fmt.Sprintf(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
❌ TRADING SIGNAL INVALIDATED
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Signal ID:  %s
Symbol:     %s
Direction:  %s
Strategy:   %s

⚠️  Signal conditions no longer met during confirmation period.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`,
		signal.SignalID,
		signal.Symbol,
		signal.Type,
		signal.StrategyName,
	)

	n.logger.Info(message)
	return nil
}

func (n *ConsoleNotifier) notifySignalOutcome(notification *Notification) error {
	signal := notification.Signal
	outcome := notification.Outcome

	if signal == nil || outcome == nil {
		return fmt.Errorf("signal or outcome is nil")
	}

	outcomeEmoji := "📊"
	if outcome.IsProfit() {
		outcomeEmoji = "💰"
	} else if outcome.IsLoss() {
		outcomeEmoji = "📉"
	}

	message := fmt.Sprintf(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
%s TRADING SIGNAL OUTCOME
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Signal ID:  %s
Symbol:     %s
Direction:  %s
Strategy:   %s
Outcome:    %s

📈 Performance:
Final Change:        %+.2f%%
Max Favorable Move:  %+.2f%%
Max Adverse Move:    %+.2f%%
Total Tracking:      %d hours

%s
%s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`,
		outcomeEmoji,
		signal.SignalID,
		signal.Symbol,
		signal.Type,
		signal.StrategyName,
		outcome.Outcome,
		outcome.FinalPriceChangePct.InexactFloat64(),
		outcome.MaxFavorableMovePct.InexactFloat64(),
		outcome.MaxAdverseMovePct.InexactFloat64(),
		outcome.TotalTrackingHours,
		conditionalField("Profit Target Hit", outcome.ProfitTargetHit),
		conditionalField("Stop Loss Hit", outcome.StopLossHit),
	)

	n.logger.Info(message)
	return nil
}

func (n *ConsoleNotifier) notifySystemError(notification *Notification) error {
	message := fmt.Sprintf(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  SYSTEM ERROR
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
%s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`,
		notification.Message,
	)

	n.logger.Error(message)
	return nil
}

func conditionalField(label string, value bool) string {
	if value {
		return fmt.Sprintf("✓ %s: YES", label)
	}
	return fmt.Sprintf("✗ %s: NO", label)
}
