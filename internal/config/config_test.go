package config_test

import (
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
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	p := writeConfig(t, `{"host":"127.0.0.1","interval":60000000000,"port_ranges":[{"from":1,"to":1024}],"state_path":"/tmp/state.json"}`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", cfg.Host)
	}
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected 60s interval, got %v", cfg.Interval)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	p := writeConfig(t, `{not valid json}`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestValidate_EmptyPortRanges(t *testing.T) {
	p := writeConfig(t, `{"port_ranges":[]}`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for empty port_ranges")
	}
}

func TestValidate_InvalidRange(t *testing.T) {
	p := writeConfig(t, `{"port_ranges":[{"from":1024,"to":80}]}`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for from > to")
	}
}

func TestValidate_WebhookDefaultTimeout(t *testing.T) {
	p := writeConfig(t, `{"port_ranges":[{"from":80,"to":80}],"webhook":{"url":"http://example.com/hook"}}`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Webhook.Timeout != 10*time.Second {
		t.Errorf("expected default 10s webhook timeout, got %v", cfg.Webhook.Timeout)
	}
}
