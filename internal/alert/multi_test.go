package alert_test

import (
	"bytes"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestMultiNotifier_Add(t *testing.T) {
	var buf bytes.Buffer
	mn := alert.NewMultiNotifier()
	mn.Add(alert.NewLogNotifier(&buf))

	a := fixedAlert(3000, alert.LevelInfo, "added dynamically")
	if err := mn.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output after dynamically added notifier, got none")
	}
}

func TestMultiNotifier_EmptyIsNoop(t *testing.T) {
	mn := alert.NewMultiNotifier()
	a := fixedAlert(9090, alert.LevelInfo, "no notifiers registered")
	if err := mn.Notify(a); err != nil {
		t.Fatalf("expected no error with zero notifiers, got: %v", err)
	}
}
