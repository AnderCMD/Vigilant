// Package vigil schedules periodic activity pulses.
package vigil

import (
	"context"
	"sync"
	"time"
)

// Pulse simulates one unit of user activity.
type Pulse func() error

// Runner calls a Pulse at a fixed interval until its context is cancelled.
// It is safe for concurrent use.
type Runner struct {
	pulse Pulse

	mu       sync.Mutex
	interval time.Duration
	paused   bool
	lastErr  error
	lastAt   time.Time
	wake     chan struct{}
}

// New creates a Runner. It does nothing until Run is called.
func New(pulse Pulse, interval time.Duration, paused bool) *Runner {
	return &Runner{pulse: pulse, interval: interval, paused: paused, wake: make(chan struct{}, 1)}
}

// Run blocks, pulsing every interval while not paused, until ctx is done.
func (r *Runner) Run(ctx context.Context) {
	timer := time.NewTimer(r.Interval())
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-r.wake:
			// Interval or pause state changed: restart the countdown.
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(r.Interval())
		case <-timer.C:
			if !r.Paused() {
				err := r.pulse()
				r.mu.Lock()
				r.lastErr, r.lastAt = err, time.Now()
				r.mu.Unlock()
			}
			timer.Reset(r.Interval())
		}
	}
}

func (r *Runner) notify() {
	select {
	case r.wake <- struct{}{}:
	default:
	}
}

// Interval returns the current pulse interval.
func (r *Runner) Interval() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.interval
}

// SetInterval changes the pulse interval and restarts the countdown.
func (r *Runner) SetInterval(d time.Duration) {
	r.mu.Lock()
	r.interval = d
	r.mu.Unlock()
	r.notify()
}

// Paused reports whether pulses are currently suppressed.
func (r *Runner) Paused() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.paused
}

// SetPaused pauses or resumes pulsing and restarts the countdown.
func (r *Runner) SetPaused(p bool) {
	r.mu.Lock()
	r.paused = p
	r.mu.Unlock()
	r.notify()
}

// Last returns the time and result of the most recent pulse.
func (r *Runner) Last() (time.Time, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lastAt, r.lastErr
}
