package monitor_test

import (
	"net"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func listenTCP(t *testing.T) (port int, close func()) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port = l.Addr().(*net.TCPAddr).Port
	return port, func() { l.Close() }
}

func newMonitor(t *testing.T) *monitor.Monitor {
	t.Helper()
	s, err := scanner.NewScanner("127.0.0.1", 500)
	if err != nil {
		t.Fatalf("NewScanner: %v", err)
	}
	return monitor.New(s)
}

func TestScan_DetectsOpenPort(t *testing.T) {
	port, closePort := listenTCP(t)
	defer closePort()

	m := newMonitor(t)
	changes, err := m.Scan(port, port)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(changes) != 1 || !changes[0].CurrOpen {
		t.Fatalf("expected one open change, got %v", changes)
	}
}

func TestScan_DetectsClose(t *testing.T) {
	port, closePort := listenTCP(t)

	m := newMonitor(t)
	// First scan — port is open
	if _, err := m.Scan(port, port); err != nil {
		t.Fatalf("first Scan: %v", err)
	}

	closePort()

	// Second scan — port should now be closed
	changes, err := m.Scan(port, port)
	if err != nil {
		t.Fatalf("second Scan: %v", err)
	}
	if len(changes) != 1 || changes[0].CurrOpen {
		t.Fatalf("expected one close change, got %v", changes)
	}
}

func TestScan_NoChangeOnRepeat(t *testing.T) {
	port, closePort := listenTCP(t)
	defer closePort()

	m := newMonitor(t)
	if _, err := m.Scan(port, port); err != nil {
		t.Fatalf("first Scan: %v", err)
	}
	changes, err := m.Scan(port, port)
	if err != nil {
		t.Fatalf("second Scan: %v", err)
	}
	if len(changes) != 0 {
		t.Fatalf("expected no changes on repeat scan, got %v", changes)
	}
}

func TestSnapshot(t *testing.T) {
	port, closePort := listenTCP(t)
	defer closePort()

	m := newMonitor(t)
	if _, err := m.Scan(port, port); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	snap := m.Snapshot()
	if !snap[port] {
		t.Fatalf("expected port %d to be open in snapshot", port)
	}
}
