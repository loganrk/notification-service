package slidingWindow

import (
	"context"
	"sync"
	"time"
)

type limiter struct {
	mu         sync.Mutex
	maxEvents  int
	window     time.Duration
	timestamps []time.Time
}

func New(maxEvents int, window time.Duration) *limiter {
	return &limiter{
		maxEvents: maxEvents,
		window:    window,
	}
}

// Allow returns true if the request can proceed immediately.
func (l *limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cleanupExpired(time.Now())

	if len(l.timestamps) < l.maxEvents {
		l.timestamps = append(l.timestamps, time.Now())
		return true
	}
	return false
}

// WaitUntilAllowed blocks until the request is allowed or the context is canceled.
func (l *limiter) WaitUntilAllowed(ctx context.Context) error {
	for {
		l.mu.Lock()
		now := time.Now()
		l.cleanupExpired(now)

		if len(l.timestamps) < l.maxEvents {
			l.timestamps = append(l.timestamps, now)
			l.mu.Unlock()
			return nil
		}

		oldest := l.timestamps[0]
		waitDuration := oldest.Add(l.window).Sub(now)
		l.mu.Unlock()

		if waitDuration <= 0 {
			waitDuration = 100 * time.Millisecond
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitDuration):
			// Retry after waiting
		}
	}
}

// cleanupExpired removes timestamps outside the sliding window.
func (l *limiter) cleanupExpired(now time.Time) {
	if len(l.timestamps) == 0 {
		return
	}

	cutoff := now.Add(-l.window)
	var filtered []time.Time
	for _, ts := range l.timestamps {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}
	l.timestamps = filtered
}
