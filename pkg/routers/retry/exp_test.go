package retry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExpRetry_RetryLoop(t *testing.T) {
	maxDelay := 10 * time.Millisecond
	ctx := context.Background()

	retry := NewExpRetry(3, 2, 2*time.Millisecond, &maxDelay)

	idx := 0
	iterator := retry.Iterator()

	for iterator.HasNext() {
		idx++

		require.NoError(t, iterator.WaitNext(ctx))
	}

	require.Equal(t, 3, idx)
}

func TestExpRetry_WaitTime(t *testing.T) {
	maxRetries := 4
	maxDelay := 10 * time.Millisecond
	expectedDelays := []time.Duration{
		2 * time.Millisecond,
		4 * time.Millisecond,
		8 * time.Millisecond,
		10 * time.Millisecond,
		10 * time.Millisecond,
	}

	retry := NewExpRetry(maxRetries, 2, 2*time.Millisecond, &maxDelay)

	iterator := retry.Iterator()

	for attempt, expectedDelay := range expectedDelays {
		require.Equal(t, expectedDelay, iterator.getNextWaitDuration(attempt))
	}
}
