// Package monitor provides port state tracking and change detection.
package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// PortState represents the open/closed state of a single port.
type PortState struct {
	Port   int
	Open   bool
	SeenAt time.Time
}

// Change describes a transition in port state.
type Change struct {
	Port     int
	PrevOpen bool
	CurrOpen bool
	Detected time.Time
}

func (c Change) String() string {
	if c.CurrOpen {
		return fmt.Sprintf("port %d opened at %s", c.Port, c.DetectedAt())
	}
	return fmt.Sprintf("port %d closed at %s", c.Port, c.DetectedAt())
}

func (c Change) DetectedAt() string {
	return c.Detected.Format(time.RFC3339)
}

// Monitor tracks port states across scans and emits changes.
type Monitor struct {
	mu       sync.Mutex
	scanner  *scanner.Scanner
	previous map[int]bool
}

// New creates a Monitor using the provided Scanner.
func New(s *scanner.Scanner) *Monitor {
	return &Monitor{
		scanner:  s,
		previous: make(map[int]bool),
	}
}

// Scan performs a port scan over [start, end] and returns any state changes.
func (m *Monitor) Scan(start, end int) ([]Change, error) {
	states, err := m.scanner.Scan(start, end)
	if err != nil {
		return nil, fmt.Errorf("monitor scan: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	var changes []Change

	for port, open := range states {
		prev, seen := m.previous[port]
		if !seen || prev != open {
			changes = append(changes, Change{
				Port:     port,
				PrevOpen: prev,
				CurrOpen: open,
				Detected: now,
			})
		}
		m.previous[port] = open
	}

	return changes, nil
}

// Snapshot returns a copy of the last known port states.
func (m *Monitor) Snapshot() map[int]bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[int]bool, len(m.previous))
	for k, v := range m.previous {
		out[k] = v
	}
	return out
}
