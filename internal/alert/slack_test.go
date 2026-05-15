package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSlackNotifier_Success(t *testing.T) {
	var received slackPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json content-type, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL)
	a := Alert{Port: 8080, Kind: "opened", Timestamp: time.Now()}

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Text == "" {
		t.Error("expected non-empty slack message text")
	}
}

func TestSlackNotifier_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewSlackNotifier(ts.URL)
	a := Alert{Port: 443, Kind: "closed", Timestamp: time.Now()}

	if err := n.Notify(a); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestSlackNotifier_UnreachableURL(t *testing.T) {
	n := NewSlackNotifier("http://127.0.0.1:0/no-server")
	a := Alert{Port: 22, Kind: "opened", Timestamp: time.Now()}

	if err := n.Notify(a); err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
