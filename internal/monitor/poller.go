package monitor

import (
	"context"
	"log"
	"time"
)

// PollConfig holds configuration for the periodic poller.
type PollConfig struct {
	StartPort int
	EndPort   int
	Interval  time.Duration
	// OnChange is called for every detected state change. It must be safe for
	// concurrent use; the poller serialises calls from a single goroutine.
	OnChange func(Change)
}

// Poll runs a blocking polling loop that scans the configured port range at
// the given interval, invoking OnChange for every detected transition.
// It returns when ctx is cancelled.
func (m *Monitor) Poll(ctx context.Context, cfg PollConfig) {
	if cfg.Interval <= 0 {
		cfg.Interval = 30 * time.Second
	}
	if cfg.OnChange == nil {
		cfg.OnChange = func(c Change) { log.Println(c) }
	}

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	// Run an immediate first scan before waiting for the first tick.
	m.runScan(cfg)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.runScan(cfg)
		}
	}
}

func (m *Monitor) runScan(cfg PollConfig) {
	changes, err := m.Scan(cfg.StartPort, cfg.EndPort)
	if err != nil {
		log.Printf("portwatch: scan error: %v", err)
		return
	}
	for _, c := range changes {
		cfg.OnChange(c)
	}
}
