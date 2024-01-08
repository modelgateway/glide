package health

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateLimitTracker_ResetCorrectly(t *testing.T) {
	tracker := NewRateLimitTracker()
	require.False(t, tracker.Limited())

	tracker.SetLimited(10 * time.Millisecond)
	require.True(t, tracker.Limited())

	time.Sleep(11 * time.Millisecond)
	require.False(t, tracker.Limited())
}
