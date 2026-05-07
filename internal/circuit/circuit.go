package circuit

import (
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

// Breaker implements a circuit breaker pattern.
type Breaker struct {
	mu          sync.Mutex
	state       State
	failures    int
	threshold   int
	timeout     time.Duration
	lastFailure time.Time
}

// New creates a circuit breaker with the given failure threshold and recovery timeout.
func New(threshold int, timeout time.Duration) *Breaker {
	return &Breaker{
		state:     Closed,
		threshold: threshold,
		timeout:   timeout,
	}
}

// Allow checks if the request should be allowed through.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case Closed:
		return true
	case Open:
		if time.Since(b.lastFailure) > b.timeout {
			b.state = HalfOpen
			return true
		}
		return false
	case HalfOpen:
		return true
	}
	return false
}

// RecordSuccess records a successful request.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = Closed
}

// RecordFailure records a failed request.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	b.lastFailure = time.Now()
	if b.failures >= b.threshold {
		b.state = Open
	}
}
