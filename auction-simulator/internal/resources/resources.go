package resources

import "time"

// Limiter controls the rate of auctions
type Limiter interface {
	Acquire()
	Release()
}

// SimpleRateLimiter implements a basic rate limiter
type SimpleRateLimiter struct {
	interval time.Duration
}

// NewSimpleRateLimiter creates a new rate limiter with the given interval
func NewSimpleRateLimiter(interval time.Duration) *SimpleRateLimiter {
	return &SimpleRateLimiter{
		interval: interval,
	}
}

// Acquire implements the Limiter interface
func (l *SimpleRateLimiter) Acquire() {
	time.Sleep(l.interval)
}

// Release implements the Limiter interface
func (l *SimpleRateLimiter) Release() {
	// No-op for this simple implementation
}

// NewLimiter creates a limiter based on available resources. The parameters
// are currently unused and kept for future tuning. It returns a Limiter
// implementation.
func NewLimiter(vcpu, ramMb int) Limiter {
	// Simple heuristic: use a small interval proportional to number of CPUs
	interval := time.Duration(10) * time.Millisecond
	if vcpu > 0 {
		interval = time.Duration(10/vcpu) * time.Millisecond
	}
	return NewSimpleRateLimiter(interval)
}

// AuctionConfig holds configuration for the auction
type AuctionConfig struct {
	NumBidders    int
	BidTimeoutSec int
	MinBidAmount  float64
	MaxBidAmount  float64
}

// DefaultConfig returns default auction configuration
func DefaultConfig() AuctionConfig {
	return AuctionConfig{
		NumBidders:    5,
		BidTimeoutSec: 3,
		MinBidAmount:  100,
		MaxBidAmount:  1000,
	}
}
