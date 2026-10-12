package vigil

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func start(t *testing.T, r *Runner) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { r.Run(ctx); close(done) }()
	t.Cleanup(func() {
		cancel()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Error("Run did not stop after cancel")
		}
	})
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("condition not met in time")
}

func TestPulsesRepeatedly(t *testing.T) {
	var n atomic.Int32
	start(t, New(func() error { n.Add(1); return nil }, 10*time.Millisecond, false))
	waitFor(t, func() bool { return n.Load() >= 3 })
}

func TestPausedDoesNotPulse(t *testing.T) {
	var n atomic.Int32
	r := New(func() error { n.Add(1); return nil }, 5*time.Millisecond, true)
	start(t, r)
	time.Sleep(60 * time.Millisecond)
	if n.Load() != 0 {
		t.Fatalf("pulsed %d times while paused", n.Load())
	}
	r.SetPaused(false)
	waitFor(t, func() bool { return n.Load() >= 1 })
}

func TestSetIntervalTakesEffect(t *testing.T) {
	var n atomic.Int32
	r := New(func() error { n.Add(1); return nil }, time.Hour, false)
	start(t, r)
	r.SetInterval(5 * time.Millisecond)
	waitFor(t, func() bool { return n.Load() >= 1 })
}

func TestRecordsLastError(t *testing.T) {
	boom := errors.New("boom")
	r := New(func() error { return boom }, 5*time.Millisecond, false)
	start(t, r)
	waitFor(t, func() bool { at, _ := r.Last(); return !at.IsZero() })
	if _, err := r.Last(); !errors.Is(err, boom) {
		t.Fatalf("got %v, want boom", err)
	}
}
