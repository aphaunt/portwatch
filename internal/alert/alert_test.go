package alert_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func fixedAlert(port int, level alert.Level, msg string) alert.Alert {
	return alert.Alert{
		Timestamp: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
		Level:     level,
		Port:      port,
		Message:   msg,
	}
}

func TestAlert_String(t *testing.T) {
	a := fixedAlert(8080, alert.LevelWarn, "port opened unexpectedly")
	s := a.String()
	if !strings.Contains(s, "WARN") {
		t.Errorf("expected WARN in output, got: %s", s)
	}
	if !strings.Contains(s, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", s)
	}
}

func TestLogNotifier_Notify(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewLogNotifier(&buf)
	a := fixedAlert(443, alert.LevelInfo, "port closed")

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "443") {
		t.Errorf("expected port 443 in log output, got: %s", buf.String())
	}
}

func TestMultiNotifier_NotifiesAll(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	n1 := alert.NewLogNotifier(&buf1)
	n2 := alert.NewLogNotifier(&buf2)
	mn := alert.NewMultiNotifier(n1, n2)

	a := fixedAlert(22, alert.LevelWarn, "ssh port opened")
	if err := mn.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, buf := range []*bytes.Buffer{&buf1, &buf2} {
		if !strings.Contains(buf.String(), "22") {
			t.Errorf("notifier %d did not receive alert", i+1)
		}
	}
}

func TestMultiNotifier_CollectsErrors(t *testing.T) {
	failing := &failNotifier{}
	mn := alert.NewMultiNotifier(failing, failing)

	err := mn.Notify(fixedAlert(80, alert.LevelError, "test"))
	if err == nil {
		t.Fatal("expected error from failing notifiers")
	}
}

type failNotifier struct{}

func (f *failNotifier) Notify(_ alert.Alert) error {
	return fmt.Errorf("simulated notifier failure")
}
