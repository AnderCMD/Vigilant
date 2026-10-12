package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadMissingReturnsDefaults(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg != Default() {
		t.Fatalf("got %+v, want defaults", cfg)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "config.json")
	want := Config{IntervalSeconds: 600, KeepDisplayAwake: false, StartPaused: true}
	if err := Save(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestLoadCorruptFallsBackToDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err == nil {
		t.Fatal("expected an error for corrupt file")
	}
	if cfg != Default() {
		t.Fatalf("got %+v, want defaults", cfg)
	}
}

func TestNormalizeClamps(t *testing.T) {
	cases := map[int]time.Duration{
		0:        DefaultInterval,
		-5:       DefaultInterval,
		1:        MinInterval,
		99999999: MaxInterval,
		120:      2 * time.Minute,
	}
	for in, want := range cases {
		c := Config{IntervalSeconds: in}
		c.Normalize()
		if c.Interval() != want {
			t.Errorf("Normalize(%d) = %v, want %v", in, c.Interval(), want)
		}
	}
}
