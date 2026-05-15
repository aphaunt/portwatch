// Package alert provides notification primitives for portwatch.
//
// An Alert describes a port state change event (opened or closed) along with
// the port number and the time the change was detected.
//
// Notifiers
//
// The package ships several Notifier implementations:
//
//   - LogNotifier  — writes alerts to a standard logger (default / fallback).
//   - WebhookNotifier — HTTP POST to an arbitrary webhook endpoint.
//   - SlackNotifier — HTTP POST to a Slack incoming-webhook URL.
//   - MultiNotifier — fan-out wrapper that dispatches to multiple notifiers
//     and collects all errors.
//
// Usage
//
//	notifier := alert.NewMultiNotifier()
//	notifier.Add(alert.NewLogNotifier(log.Default()))
//	notifier.Add(alert.NewSlackNotifier(os.Getenv("SLACK_WEBHOOK_URL")))
//
//	// later, on a detected change:
//	notifier.Notify(alert.Alert{Port: 8080, Kind: "opened", Timestamp: time.Now()})
package alert
