// Package config handles loading and validating portwatch configuration.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// PortRange describes an inclusive range of TCP ports to monitor.
type PortRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// WebhookConfig holds optional webhook notifier settings.
type WebhookConfig struct {
	URL     string        `json:"url"`
	Timeout time.Duration `json:"timeout"`
}

// Config is the top-level portwatch configuration structure.
type Config struct {
	Host       string        `json:"host"`
	Interval   time.Duration `json:"interval"`
	PortRanges []PortRange   `json:"port_ranges"`
	StatePath  string        `json:"state_path"`
	Webhook    *WebhookConfig `json:"webhook,omitempty"`
}

// Load reads and unmarshals a JSON config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %s: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Validate checks that the configuration is semantically correct.
func (c *Config) Validate() error {
	if len(c.PortRanges) == 0 {
		return errors.New("config: at least one port_range is required")
	}
	for i, pr := range c.PortRanges {
		if pr.From < 1 || pr.To > 65535 {
			return fmt.Errorf("config: port_range[%d] out of valid range 1-65535", i)
		}
		if pr.From > pr.To {
			return fmt.Errorf("config: port_range[%d] from (%d) must be <= to (%d)", i, pr.From, pr.To)
		}
	}
	if c.Interval <= 0 {
		c.Interval = 30 * time.Second
	}
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Webhook != nil && c.Webhook.Timeout <= 0 {
		c.Webhook.Timeout = 10 * time.Second
	}
	return nil
}
