package leakyBucket

import (
	"context"
	"sync"
	"time"
)

type limiter struct {
	mu          sync.Mutex
	capacity    int           // max tokens in bucket
	interval    time.Duration // total time window (e.g., 1s for 10 reqs/sec)
	tokens      int           // current token count
	lastLeak    time.Time     // last time tokens were leaked
	leakSpacing time.Duration // time between token refill (interval / capacity)
}

// New initializes a leaky bucket limiter.
func New(capacity int, interval time.Duration) *limiter {
	return &limiter{
		capacity:    capacity,
		interval:    interval,
		tokens:      capacity,
		lastLeak:    time.Now(),
		leakSpacing: interval / time.Duration(capacity),
	}
}

// Allow returns true if the request can be served immediately.
func (l *limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.refill(now)

	if l.tokens > 0 {
		l.tokens--
		return true
	}

	return false
}

// WaitUntilAllowed blocks until the request is allowed or context is cancelled.
func (l *limiter) WaitUntilAllowed(ctx context.Context) error {
	for {
		l.mu.Lock()
		now := time.Now()
		l.refill(now)

		if l.tokens > 0 {
			l.tokens--
			l.mu.Unlock()
			return nil
		}

		wait := l.leakSpacing - now.Sub(l.lastLeak)
		if wait < 0 {
			wait = l.leakSpacing
		}
		l.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
			// Retry after wait
		}
	}
}

// refill adds tokens based on time elapsed since last leak.
func (l *limiter) refill(now time.Time) {
	elapsed := now.Sub(l.lastLeak)
	newTokens := int(elapsed / l.leakSpacing)
	if newTokens > 0 {
		l.tokens = min(l.capacity, l.tokens+newTokens)
		l.lastLeak = l.lastLeak.Add(time.Duration(newTokens) * l.leakSpacing)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
