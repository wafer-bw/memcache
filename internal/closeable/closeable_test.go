package closeable_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/internal/closeable"
)

func TestClose(t *testing.T) {
	t.Parallel()

	t.Run("subsequent calls do not panic", func(t *testing.T) {
		t.Parallel()

		c := closeable.New()
		c.Close()
		c.Close()
	})
}

func TestClosed(t *testing.T) {
	t.Parallel()

	t.Run("returns false when not closed", func(t *testing.T) {
		t.Parallel()

		c := closeable.New()
		close := c.Closed()
		require.False(t, close)
	})

	t.Run("returns true when closed", func(t *testing.T) {
		t.Parallel()

		c := closeable.New()
		c.Close()
		closed := c.Closed()
		require.True(t, closed)
	})
}
