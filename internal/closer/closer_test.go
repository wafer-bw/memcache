package closer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/closer"
)

func TestClose(t *testing.T) {
	t.Parallel()

	t.Run("subsequent calls do not panic", func(t *testing.T) {
		t.Parallel()

		c := closer.New()
		c.Close()
		c.Close()
	})
}

func TestWaitClosed(t *testing.T) {
	t.Parallel()

	t.Run("returns channel for waiting on when not closed", func(t *testing.T) {
		t.Parallel()

		c := closer.New()
		select {
		case <-c.WaitClosed():
			t.Fatal("should not have closed")
		default:
			// pass
		}
	})

	t.Run("returns closed channel when closed", func(t *testing.T) {
		t.Parallel()

		c := closer.New()
		c.Close()
		select {
		case <-c.WaitClosed():
			// pass
		default:
			t.Fatal("should have closed")
		}
	})
}

func TestClosed(t *testing.T) {
	t.Parallel()

	t.Run("returns false when not closed", func(t *testing.T) {
		t.Parallel()

		c := closer.New()
		close := c.Closed()
		require.False(t, close)
	})

	t.Run("returns true when closed", func(t *testing.T) {
		t.Parallel()

		c := closer.New()
		c.Close()
		closed := c.Closed()
		require.True(t, closed)
	})
}
