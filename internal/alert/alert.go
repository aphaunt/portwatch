// Package alert defines the Notifier interface and built-in implementations
// for delivering port-change alerts to various destinations.
package alert

import (
	"fmt"
	"log"
	"time"
)

// EventType describes the kind of port change that triggered an alert.
type EventType int

const (
	EventOpened EventType = iota
	EventClosed
)

// String returns a human-readable label for the event.
func (e EventType) String() string {
	switch e {
	case EventOpened:
		return "opened"
	case EventClosed:
		return "closed"
	default:
		return "unknown"
	}
}

// Alert carries the details of a single port-change event.
type Alert struct {
	Timestamp time.Time
	Event     EventType
	Port      int
	Host      string
}

// String returns a formatted one-line description of the alert.
func (a Alert) String() string {
	return fmt.Sprintf("[%s] port %d %s on %s",
		a.Timestamp.UTC().Format(time.RFC3339), a.Port, a.Event, a.Host)
}

// Notifier is the interface implemented by all alert back-ends.
type Notifier interface {
	Notify(Alert) error
}

// LogNotifier writes alerts to the standard logger.
type LogNotifier struct {
	logger *log.Logger
}

// NewLogNotifier creates a LogNotifier backed by the provided logger.
// If logger is nil, the default log package logger is used.
func NewLogNotifier(logger *log.Logger) *LogNotifier {
	if logger == nil {
		logger = log.Default()
	}
	return &LogNotifier{logger: logger}
}

// Notify logs the alert and always returns nil.
func (l *LogNotifier) Notify(a Alert) error {
	l.logger.Println(a.String())
	return nil
}
