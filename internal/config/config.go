// Package config loads and stores the user's Vigilant settings.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// MinInterval is the shortest allowed pulse interval.
	MinInterval = 10 * time.Second
	// MaxInterval is the longest allowed pulse interval.
	MaxInterval = time.Hour
	// DefaultInterval is used when no (valid) interval is configured.
	DefaultInterval = 3 * time.Minute
)

// Config holds the persisted user settings.
type Config struct {
	// IntervalSeconds is the time between two activity pulses.
	IntervalSeconds int `json:"intervalSeconds"`
	// KeepDisplayAwake also asks Windows to keep the screen on.
	KeepDisplayAwake bool `json:"keepDisplayAwake"`
	// StartPaused launches Vigilant in the paused state.
	StartPaused bool `json:"startPaused"`
}

// Default returns the default configuration.
func Default() Config {
	return Config{
		IntervalSeconds:  int(DefaultInterval / time.Second),
		KeepDisplayAwake: true,
	}
}

// Interval returns the pulse interval as a duration.
func (c Config) Interval() time.Duration {
	return time.Duration(c.IntervalSeconds) * time.Second
}

// Normalize clamps out-of-range values back into the supported range.
func (c *Config) Normalize() {
	d := c.Interval()
	switch {
	case d <= 0:
		c.IntervalSeconds = int(DefaultInterval / time.Second)
	case d < MinInterval:
		c.IntervalSeconds = int(MinInterval / time.Second)
	case d > MaxInterval:
		c.IntervalSeconds = int(MaxInterval / time.Second)
	}
}

// Dir returns the directory where Vigilant stores its settings.
func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "Vigilant"), nil
}

// Path returns the full path of the configuration file.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load reads the configuration from path. A missing file yields the defaults;
// a corrupt file also yields the defaults together with an error describing it.
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), fmt.Errorf("parse %s: %w", path, err)
	}
	cfg.Normalize()
	return cfg, nil
}

// Save writes the configuration to path atomically.
func Save(path string, cfg Config) error {
	cfg.Normalize()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
