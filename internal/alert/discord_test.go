package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDiscordNotifier_Success(t *testing.T) {
	var received discordPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL)
	a := Alert{Port: 8080, Kind: "opened", Host: "localhost"}

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(received.Content, "portwatch alert") {
		t.Errorf("expected content to contain 'portwatch alert', got: %s", received.Content)
	}
	if !strings.Contains(received.Content, a.String()) {
		t.Errorf("expected content to contain alert string, got: %s", received.Content)
	}
}

func TestDiscordNotifier_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer ts.Close()

	n := NewDiscordNotifier(ts.URL)
	a := Alert{Port: 443, Kind: "closed", Host: "localhost"}

	err := n.Notify(a)
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
	if !strings.Contains(err.Error(), "429") {
		t.Errorf("expected error to mention status 429, got: %v", err)
	}
}

func TestDiscordNotifier_UnreachableURL(t *testing.T) {
	n := NewDiscordNotifier("http://127.0.0.1:1")
	a := Alert{Port: 22, Kind: "opened", Host: "localhost"}

	err := n.Notify(a)
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
