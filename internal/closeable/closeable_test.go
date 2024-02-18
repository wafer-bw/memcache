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
		require.False(t, c.Closed())
	})

	t.Run("returns true when closed", func(t *testing.T) {
		t.Parallel()

		c := closeable.New()
		c.Close()
		require.True(t, c.Closed())
	})
}
