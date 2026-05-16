package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestOpsGenieNotifier_Success(t *testing.T) {
	var received ogPayload
	var authHeader string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n := NewOpsGenieNotifier("og-test-key")
	n.alertURL = ts.URL

	a := Alert{Port: 9090, Kind: KindOpened, Time: time.Now()}
	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(authHeader, "og-test-key") {
		t.Errorf("auth header = %q, want GenieKey og-test-key", authHeader)
	}
	if received.Source != "portwatch" {
		t.Errorf("source = %q, want portwatch", received.Source)
	}
	if received.Message == "" {
		t.Error("message should not be empty")
	}
	if len(received.Tags) == 0 {
		t.Error("tags should not be empty")
	}
}

func TestOpsGenieNotifier_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n := NewOpsGenieNotifier("bad-key")
	n.alertURL = ts.URL

	err := n.Notify(Alert{Port: 80, Kind: KindClosed, Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestOpsGenieNotifier_UnreachableURL(t *testing.T) {
	n := NewOpsGenieNotifier("key")
	n.alertURL = "http://127.0.0.1:1"

	err := n.Notify(Alert{Port: 443, Kind: KindOpened, Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
