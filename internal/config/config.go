// Package config provides configuration loading and validation for portwatch.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Config holds the top-level portwatch configuration.
type Config struct {
	// PortRanges is a list of port ranges to monitor, e.g. ["1-1024", "8080", "9000-9100"].
	PortRanges []string `json:"port_ranges"`

	// Interval is how often to scan ports.
	Interval Duration `json:"interval"`

	// LogPath is an optional file path to write alerts. Empty means stdout.
	LogPath string `json:"log_path,omitempty"`
}

// Duration is a wrapper around time.Duration that supports JSON unmarshalling
// from a human-readable string such as "30s" or "1m".
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("config: invalid duration %q: %w", s, err)
	}
	d.Duration = v
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

// Load reads and parses a JSON config file at the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Validate checks that required fields are present and values are sensible.
func (c *Config) Validate() error {
	if len(c.PortRanges) == 0 {
		return errors.New("config: port_ranges must not be empty")
	}
	if c.Interval.Duration <= 0 {
		return errors.New("config: interval must be a positive duration")
	}
	return nil
}
