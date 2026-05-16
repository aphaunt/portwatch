package alert

import (
	"fmt"
	"log/syslog"
)

// SyslogNotifier sends alerts to the local syslog daemon.
type SyslogNotifier struct {
	writer *syslog.Writer
	tag    string
}

// NewSyslogNotifier creates a SyslogNotifier that writes to syslog with the
// given tag (e.g. "portwatch"). Priority is LOG_WARNING | LOG_DAEMON.
func NewSyslogNotifier(tag string) (*SyslogNotifier, error) {
	w, err := syslog.New(syslog.LOG_WARNING|syslog.LOG_DAEMON, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog: open writer: %w", err)
	}
	return &SyslogNotifier{writer: w, tag: tag}, nil
}

// Notify sends the alert to syslog.
func (s *SyslogNotifier) Notify(a Alert) error {
	msg := fmt.Sprintf("portwatch alert: %s", a.String())
	if err := s.writer.Warning(msg); err != nil {
		return fmt.Errorf("syslog: write warning: %w", err)
	}
	return nil
}

// Close releases the underlying syslog connection.
func (s *SyslogNotifier) Close() error {
	if err := s.writer.Close(); err != nil {
		return fmt.Errorf("syslog: close writer: %w", err)
	}
	return nil
}
