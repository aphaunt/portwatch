package alert

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// TestMultiNotifier_PagerDutyAndOpsGenie verifies that both PagerDuty and
// OpsGenie notifiers can be combined via MultiNotifier and both receive events.
func TestMultiNotifier_PagerDutyAndOpsGenie(t *testing.T) {
	var pdCalls, ogCalls int32

	pdServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&pdCalls, 1)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer pdServer.Close()

	ogServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&ogCalls, 1)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ogServer.Close()

	pd := NewPagerDutyNotifier("pd-key")
	pd.eventURL = pdServer.URL

	og := NewOpsGenieNotifier("og-key")
	og.alertURL = ogServer.URL

	multi := NewMultiNotifier(pd, og)

	a := Alert{Port: 8443, Kind: KindOpened, Time: time.Now()}
	if err := multi.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if atomic.LoadInt32(&pdCalls) != 1 {
		t.Errorf("pagerduty calls = %d, want 1", pdCalls)
	}
	if atomic.LoadInt32(&ogCalls) != 1 {
		t.Errorf("opsgenie calls = %d, want 1", ogCalls)
	}
}

// TestMultiNotifier_ContinuesOnPartialFailure verifies that a failure in one
// notifier does not prevent others from being called.
func TestMultiNotifier_ContinuesOnPartialFailure(t *testing.T) {
	var ogCalls int32

	ogServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&ogCalls, 1)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ogServer.Close()

	pd := NewPagerDutyNotifier("key")
	pd.eventURL = "http://127.0.0.1:1" // unreachable

	og := NewOpsGenieNotifier("og-key")
	og.alertURL = ogServer.URL

	multi := NewMultiNotifier(pd, og)

	a := Alert{Port: 22, Kind: KindClosed, Time: time.Now()}
	err := multi.Notify(a)
	if err == nil {
		t.Fatal("expected combined error from failing notifier")
	}
	if atomic.LoadInt32(&ogCalls) != 1 {
		t.Errorf("opsgenie should still be called; calls = %d", ogCalls)
	}
}
