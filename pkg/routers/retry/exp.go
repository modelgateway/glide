package retry

import (
	"context"
	"time"
)

// ExpRetry increase wait time exponentially with try number (delay = minDelay * baseMultiplier ^ attempt)
type ExpRetry struct {
	maxRetries     int
	baseMultiplier int
	minDelay       time.Duration
	maxDelay       *time.Duration
}

func NewExpRetry(maxRetries int, baseMultiplier int, minDelay time.Duration, maxDelay *time.Duration) *ExpRetry {
	return &ExpRetry{
		maxRetries:     maxRetries,
		baseMultiplier: baseMultiplier,
		minDelay:       minDelay,
		maxDelay:       maxDelay,
	}
}

func (r *ExpRetry) Iterator() *ExpRetryIterator {
	return &ExpRetryIterator{
		attempt:        0,
		maxRetries:     r.maxRetries,
		baseMultiplier: r.baseMultiplier,
		minDelay:       r.minDelay,
		maxDelay:       r.maxDelay,
	}
}

type ExpRetryIterator struct {
	attempt        int
	maxRetries     int
	baseMultiplier int
	minDelay       time.Duration
	maxDelay       *time.Duration
}

func (i *ExpRetryIterator) HasNext() bool {
	return i.attempt < i.maxRetries
}

func (i *ExpRetryIterator) getNextWaitDuration(attempt int) time.Duration {
	delay := i.minDelay

	if attempt > 0 {
		delay = time.Duration(float64(delay) * float64(i.baseMultiplier<<(attempt-1)))
	}

	if delay < i.minDelay {
		delay = i.minDelay
	}

	if i.maxDelay != nil && delay > *i.maxDelay {
		delay = *i.maxDelay
	}

	return delay
}

func (i *ExpRetryIterator) WaitNext(ctx context.Context) error {
	t := time.NewTimer(i.getNextWaitDuration(i.attempt))
	i.attempt++

	defer t.Stop()

	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
