package monitor_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

func TestPoll_EmitsChangeOnOpen(t *testing.T) {
	port, closePort := listenTCP(t)
	defer closePort()

	m := newMonitor(t)

	var mu sync.Mutex
	var received []monitor.Change

	cfg := monitor.PollConfig{
		StartPort: port,
		EndPort:   port,
		Interval:  20 * time.Millisecond,
		OnChange: func(c monitor.Change) {
			mu.Lock()
			received = append(received, c)
			mu.Unlock()
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go m.Poll(ctx, cfg)
	<-ctx.Done()

	mu.Lock()
	defer mu.Unlock()

	if len(received) == 0 {
		t.Fatal("expected at least one change event from poller")
	}
	if !received[0].CurrOpen {
		t.Fatalf("first change should be open, got %+v", received[0])
	}
}

func TestPoll_StopsOnContextCancel(t *testing.T) {
	port, closePort := listenTCP(t)
	defer closePort()

	m := newMonitor(t)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		m.Poll(ctx, monitor.PollConfig{
			StartPort: port,
			EndPort:   port,
			Interval:  10 * time.Millisecond,
		})
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Poll did not stop after context cancellation")
	}
}
