// Package config handles loading and validating portwatch configuration.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// NotifierConfig holds configuration for a single alert notifier.
type NotifierConfig struct {
	Type    string            `json:"type"`    // e.g. "slack", "webhook", "email", "pagerduty", "opsgenie", "discord", "teams", "syslog"
	Options map[string]string `json:"options"` // notifier-specific key/value options
}

// PortRange defines an inclusive range of ports to monitor.
type PortRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// Config is the top-level portwatch configuration structure.
type Config struct {
	// PortRanges is the list of port ranges to scan.
	PortRanges []PortRange `json:"port_ranges"`

	// Interval is how often the scanner polls, e.g. "30s", "1m".
	Interval string `json:"interval"`

	// StateFile is the path where port state is persisted between runs.
	StateFile string `json:"state_file"`

	// Notifiers is the list of alert notifier configurations.
	Notifiers []NotifierConfig `json:"notifiers"`

	// parsed holds the parsed interval duration (populated by Validate).
	parsed time.Duration
}

// IntervalDuration returns the parsed polling interval as a time.Duration.
// Validate must be called before using this method.
func (c *Config) IntervalDuration() time.Duration {
	return c.parsed
}

// Load reads and parses the JSON config file at the given path.
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
		return nil, fmt.Errorf("parsing config JSON: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Validate checks that all required fields are present and valid.
// It also parses the Interval string into a time.Duration.
func (c *Config) Validate() error {
	if len(c.PortRanges) == 0 {
		return errors.New("port_ranges must not be empty")
	}

	for i, pr := range c.PortRanges {
		if pr.From < 1 || pr.From > 65535 {
			return fmt.Errorf("port_ranges[%d].from must be between 1 and 65535", i)
		}
		if pr.To < 1 || pr.To > 65535 {
			return fmt.Errorf("port_ranges[%d].to must be between 1 and 65535", i)
		}
		if pr.From > pr.To {
			return fmt.Errorf("port_ranges[%d].from (%d) must be <= to (%d)", i, pr.From, pr.To)
		}
	}

	if c.Interval == "" {
		return errors.New("interval must not be empty")
	}

	d, err := time.ParseDuration(c.Interval)
	if err != nil {
		return fmt.Errorf("interval %q is not a valid duration: %w", c.Interval, err)
	}
	if d <= 0 {
		return errors.New("interval must be a positive duration")
	}
	c.parsed = d

	if c.StateFile == "" {
		return errors.New("state_file must not be empty")
	}

	return nil
}
