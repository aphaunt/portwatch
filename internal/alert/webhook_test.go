package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func TestWebhookNotifier_Success(t *testing.T) {
	var received map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := alert.NewWebhookNotifier(srv.URL, 5*time.Second)
	a := alert.Alert{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Event:     alert.EventOpened,
		Port:      8080,
		Host:      "localhost",
	}

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["port"] != float64(8080) {
		t.Errorf("expected port 8080, got %v", received["port"])
	}
	if received["event"] != "opened" {
		t.Errorf("expected event 'opened', got %v", received["event"])
	}
}

func TestWebhookNotifier_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := alert.NewWebhookNotifier(srv.URL, 5*time.Second)
	a := alert.Alert{Timestamp: time.Now(), Event: alert.EventClosed, Port: 22, Host: "localhost"}

	if err := n.Notify(a); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestWebhookNotifier_UnreachableURL(t *testing.T) {
	n := alert.NewWebhookNotifier("http://127.0.0.1:1", 500*time.Millisecond)
	a := alert.Alert{Timestamp: time.Now(), Event: alert.EventOpened, Port: 9999, Host: "localhost"}

	if err := n.Notify(a); err == nil {
		t.Fatal("expected connection error, got nil")
	}
}
