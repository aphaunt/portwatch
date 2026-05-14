package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeConfig: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	p := writeConfig(t, `{"port_ranges":["1-1024","8080"],"interval":"30s"}`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.PortRanges) != 2 {
		t.Errorf("expected 2 port ranges, got %d", len(cfg.PortRanges))
	}
	if cfg.Interval.Duration != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval.Duration)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	p := writeConfig(t, `not json`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestValidate_EmptyPortRanges(t *testing.T) {
	cfg := &config.Config{
		PortRanges: nil,
		Interval:   config.Duration{Duration: 10 * time.Second},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for empty port_ranges")
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := &config.Config{
		PortRanges: []string{"80"},
		Interval:   config.Duration{Duration: 0},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for zero interval")
	}
}

func TestDuration_RoundTrip(t *testing.T) {
	orig := config.Duration{Duration: 5 * time.Minute}
	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got config.Duration
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Duration != orig.Duration {
		t.Errorf("round-trip mismatch: got %v, want %v", got.Duration, orig.Duration)
	}
}
