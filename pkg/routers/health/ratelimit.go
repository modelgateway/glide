package health

import "time"

// RateLimitTracker handles rate/quota limits that often represented via 429 errors and
// has some well-defined cooldown period
type RateLimitTracker struct {
	resetAt *time.Time
}

func NewRateLimitTracker() *RateLimitTracker {
	return &RateLimitTracker{
		resetAt: nil,
	}
}

func (t *RateLimitTracker) Limited() bool {
	if t.resetAt != nil && time.Now().After(*t.resetAt) {
		t.resetAt = nil
	}

	return t.resetAt != nil
}

func (t *RateLimitTracker) SetLimited(untilReset time.Duration) {
	resetAt := time.Now().Add(untilReset)

	t.resetAt = &resetAt
}
