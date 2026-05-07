package ratelimit

import (
	"sync"
	"time"
)

// Limiter implements a token bucket rate limiter.
type Limiter struct {
	mu       sync.Mutex
	rate     float64
	burst    int
	tokens   float64
	lastTime time.Time
}

// New creates a rate limiter with the given rate (tokens/sec) and burst size.
func New(rate float64, burst int) *Limiter {
	return &Limiter{
		rate:     rate,
		burst:    burst,
		tokens:   float64(burst),
		lastTime: time.Now(),
	}
}

// Allow reports whether an event may happen now.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastTime).Seconds()
	l.lastTime = now

	l.tokens += elapsed * l.rate
	if l.tokens > float64(l.burst) {
		l.tokens = float64(l.burst)
	}

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}
