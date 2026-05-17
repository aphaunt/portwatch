// Package config handles loading and validating portwatch configuration.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// NotifierConfig holds configuration for a single alert notifier.
type NotifierConfig struct {
	Type    string            `json:"type"`
	Options map[string]string `json:"options"`
}

// PortRange defines an inclusive range of ports to monitor.
type PortRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// Config is the top-level portwatch configuration structure.
type Config struct {
	// PortRanges defines which port ranges to scan.
	PortRanges []PortRange `json:"port_ranges"`

	// IntervalSeconds is how often (in seconds) to poll ports.
	IntervalSeconds int `json:"interval_seconds"`

	// StateFile is the path where port state is persisted between runs.
	StateFile string `json:"state_file"`

	// Notifiers lists the alert backends to notify on changes.
	Notifiers []NotifierConfig `json:"notifiers"`

	// Host is the target host to scan (defaults to localhost).
	Host string `json:"host"`
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Apply defaults.
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.IntervalSeconds <= 0 {
		cfg.IntervalSeconds = 60
	}
	if cfg.StateFile == "" {
		cfg.StateFile = "/var/lib/portwatch/state.json"
	}

	return &cfg, nil
}

// Validate checks that the configuration is semantically valid.
func (c *Config) Validate() error {
	if len(c.PortRanges) == 0 {
		return errors.New("port_ranges must not be empty")
	}
	for i, r := range c.PortRanges {
		if r.From < 1 || r.From > 65535 {
			return fmt.Errorf("port_ranges[%d].from %d is out of valid range (1-65535)", i, r.From)
		}
		if r.To < 1 || r.To > 65535 {
			return fmt.Errorf("port_ranges[%d].to %d is out of valid range (1-65535)", i, r.To)
		}
		if r.From > r.To {
			return fmt.Errorf("port_ranges[%d]: from (%d) must be <= to (%d)", i, r.From, r.To)
		}
	}
	return nil
}
