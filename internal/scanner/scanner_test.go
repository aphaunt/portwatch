package scanner

import (
	"net"
	"testing"
	"time"
)

// startTestServer opens a TCP listener on a random port and returns the port and a stop func.
func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestNewScanner(t *testing.T) {
	s := NewScanner("127.0.0.1", time.Second)
	if s.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", s.Host)
	}
}

func TestScan_OpenPort(t *testing.T) {
	port, stop := startTestServer(t)
	defer stop()

	s := NewScanner("127.0.0.1", time.Second)
	results, err := s.Scan("tcp", port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScan_ClosedPort(t *testing.T) {
	s := NewScanner("127.0.0.1", 200*time.Millisecond)
	// Port 1 is almost certainly closed in test environments.
	results, err := s.Scan("tcp", 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestScan_InvalidRange(t *testing.T) {
	s := NewScanner("127.0.0.1", time.Second)
	_, err := s.Scan("tcp", 100, 50)
	if err == nil {
		t.Error("expected error for invalid port range, got nil")
	}
}

func TestScan_ResultCount(t *testing.T) {
	s := NewScanner("127.0.0.1", 100*time.Millisecond)
	results, err := s.Scan("tcp", 1, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 5 {
		t.Errorf("expected 5 results, got %d", len(results))
	}
}
