// Package alert provides alerting mechanisms for portwatch.
// It defines the Alert type and various notifier implementations
// that can be used to report unexpected port state changes.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Alert represents a notification about a port state change.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Port      int
	Message   string
}

// String returns a human-readable representation of the alert.
func (a Alert) String() string {
	return fmt.Sprintf("[%s] %s port=%d msg=%q",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Port,
		a.Message,
	)
}

// Notifier is the interface implemented by alert sinks.
type Notifier interface {
	Notify(a Alert) error
}

// LogNotifier writes alerts as text lines to an io.Writer.
type LogNotifier struct {
	Out io.Writer
}

// NewLogNotifier returns a LogNotifier that writes to w.
// If w is nil, os.Stderr is used.
func NewLogNotifier(w io.Writer) *LogNotifier {
	if w == nil {
		w = os.Stderr
	}
	return &LogNotifier{Out: w}
}

// Notify writes the alert to the underlying writer.
func (l *LogNotifier) Notify(a Alert) error {
	_, err := fmt.Fprintln(l.Out, a.String())
	return err
}
