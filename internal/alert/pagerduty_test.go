package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPagerDutyNotifier_Success(t *testing.T) {
	var received pdPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("test-key")
	n.eventURL = ts.URL

	a := Alert{Port: 8080, Kind: KindOpened, Time: time.Now()}
	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.RoutingKey != "test-key" {
		t.Errorf("routing key = %q, want %q", received.RoutingKey, "test-key")
	}
	if received.EventAction != "trigger" {
		t.Errorf("event_action = %q, want trigger", received.EventAction)
	}
	if received.Payload.Source != "portwatch" {
		t.Errorf("source = %q, want portwatch", received.Payload.Source)
	}
	if received.Payload.Summary == "" {
		t.Error("summary should not be empty")
	}
}

func TestPagerDutyNotifier_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	n := NewPagerDutyNotifier("key")
	n.eventURL = ts.URL

	err := n.Notify(Alert{Port: 443, Kind: KindClosed, Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestPagerDutyNotifier_UnreachableURL(t *testing.T) {
	n := NewPagerDutyNotifier("key")
	n.eventURL = "http://127.0.0.1:1" // nothing listening

	err := n.Notify(Alert{Port: 22, Kind: KindOpened, Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
