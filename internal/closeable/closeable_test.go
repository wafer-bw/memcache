package closeable_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/closeable"
)

func TestCloser_Ch(t *testing.T) {
	t.Parallel()

	t.Run("returns channel for waiting on when not closed", func(t *testing.T) {
		t.Parallel()

		c := closeable.New()
		select {
		case <-c.Ch():
			t.Fatal("should not have closed")
		default:
			// pass
		}
	})

	t.Run("returns closed channel when closed", func(t *testing.T) {
		t.Parallel()

		c := closeable.New()
		c.Close()
		select {
		case <-c.Ch():
			// pass
		default:
			t.Fatal("should have closed")
		}
	})
}

func TestCloser_Close(t *testing.T) {
	t.Parallel()

	t.Run("subsequent calls do not panic", func(t *testing.T) {
		t.Parallel()

		c := closeable.New()
		c.Close()
		c.Close()
	})
}

func TestCloser_Closed(t *testing.T) {
	t.Parallel()

	t.Run("returns false when not closed", func(t *testing.T) {
		t.Parallel()

		c := closeable.New()
		require.False(t, c.Closed())
	})

	t.Run("returns true when closed", func(t *testing.T) {
		t.Parallel()

		c := closeable.New()
		c.Close()
		require.True(t, c.Closed())
	})
}
