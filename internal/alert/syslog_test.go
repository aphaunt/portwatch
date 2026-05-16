package alert_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// TestSyslogNotifier_New verifies that a SyslogNotifier can be constructed
// and closed without error on platforms that expose /dev/log or equivalent.
func TestSyslogNotifier_New(t *testing.T) {
	n, err := alert.NewSyslogNotifier("portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable on this platform: %v", err)
	}
	t.Cleanup(func() { _ = n.Close() })
}

// TestSyslogNotifier_Notify sends a real alert through syslog and expects no
// error. The test is skipped when syslog is not reachable (e.g. CI containers).
func TestSyslogNotifier_Notify(t *testing.T) {
	n, err := alert.NewSyslogNotifier("portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable on this platform: %v", err)
	}
	t.Cleanup(func() { _ = n.Close() })

	a := alert.Alert{
		Port:      8080,
		Status:    alert.StatusOpened,
		Timestamp: time.Now(),
	}

	if err := n.Notify(a); err != nil {
		t.Fatalf("Notify() unexpected error: %v", err)
	}
}

// TestSyslogNotifier_Close verifies that Close is idempotent in the sense
// that calling it once does not panic.
func TestSyslogNotifier_Close(t *testing.T) {
	n, err := alert.NewSyslogNotifier("portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable on this platform: %v", err)
	}

	if err := n.Close(); err != nil {
		t.Fatalf("Close() unexpected error: %v", err)
	}
}
