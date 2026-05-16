package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTeamsNotifier_Success(t *testing.T) {
	var received teamsPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected Content-Type: %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL)
	a := Alert{Port: 8080, Kind: "opened", Timestamp: time.Now()}

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Type != "message" {
		t.Errorf("expected type 'message', got %q", received.Type)
	}
	if len(received.Attachments) == 0 {
		t.Fatal("expected at least one attachment")
	}
	blocks := received.Attachments[0].Content.Body
	if len(blocks) < 2 {
		t.Fatalf("expected at least 2 body blocks, got %d", len(blocks))
	}
	if blocks[1].Text != a.String() {
		t.Errorf("expected body text %q, got %q", a.String(), blocks[1].Text)
	}
}

func TestTeamsNotifier_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	n := NewTeamsNotifier(ts.URL)
	a := Alert{Port: 443, Kind: "closed", Timestamp: time.Now()}

	err := n.Notify(a)
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestTeamsNotifier_UnreachableURL(t *testing.T) {
	n := NewTeamsNotifier("http://127.0.0.1:0/webhook")
	a := Alert{Port: 22, Kind: "opened", Timestamp: time.Now()}

	err := n.Notify(a)
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
