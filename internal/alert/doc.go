// Package alert defines the Alert type and the Notifier interface used by
// portwatch to report unexpected port state changes.
//
// # Alert
//
// An Alert carries a timestamp, severity level, affected port number, and a
// human-readable message describing the change.
//
// # Notifier
//
// Notifier is a simple one-method interface:
//
//	type Notifier interface {
//		Notify(a Alert) error
//	}
//
// Built-in implementations:
//   - LogNotifier  — writes formatted alerts to any io.Writer (default: stderr).
//   - MultiNotifier — fans out to multiple Notifiers, collecting all errors.
//
// Additional notifiers (e.g. webhook, email) can be added by implementing the
// Notifier interface.
package alert
