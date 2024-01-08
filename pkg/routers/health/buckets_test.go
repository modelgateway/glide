package health

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTokenBucket_Take(t *testing.T) {
	bucketSize := 10
	bucket := NewTokenBucket(1000, uint64(bucketSize))

	for i := 0; i < bucketSize-1; i++ {
		require.NoError(t, bucket.Take(1))
		require.True(t, bucket.HasTokens())
	}

	// consuming 10th token
	require.NoError(t, bucket.Take(1))

	// only 10 tokens in the bucket
	require.ErrorIs(t, bucket.Take(1), ErrNoTokens)
	require.False(t, bucket.HasTokens())
}

func TestTokenBucket_TakeConcurrently(t *testing.T) {
	bucket := NewTokenBucket(10_000, 1)
	wg := &sync.WaitGroup{}

	before := time.Now()

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for k := 0; k < 10; k++ {
				for bucket.Take(1) != nil {
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()
	}

	wg.Wait()

	if time.Since(before) < 1*time.Second {
		t.Fatal("Did not wait 1s")
	}
}

func TestTokenBucket_TokenNumberIsCorrect(t *testing.T) {
	bucket := NewTokenBucket(1_000_000, 10)
	require.InEpsilon(t, 10.0, bucket.Tokens(), 0.0001)

	require.NoError(t, bucket.Take(2))
	require.InEpsilon(t, 8.0, bucket.Tokens(), 0.0001)

	require.NoError(t, bucket.Take(2))
	require.InEpsilon(t, 6.0, bucket.Tokens(), 0.0001)

	require.NoError(t, bucket.Take(2))
	require.InEpsilon(t, 4.0, bucket.Tokens(), 0.0001)

	require.NoError(t, bucket.Take(2))
	require.InEpsilon(t, 2.0, bucket.Tokens(), 0.0001)

	require.NoError(t, bucket.Take(2))
	require.LessOrEqual(t, 0.0, bucket.Tokens())
}

func TestTokenBucket_TakeBurstly(t *testing.T) {
	bucket := NewTokenBucket(1, 10)

	require.NoError(t, bucket.Take(10))
}
