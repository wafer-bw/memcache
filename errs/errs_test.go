package errs_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/memcache/errs"
)

func TestInvalidCapacityError_Error(t *testing.T) {
	t.Parallel()

	t.Run("returns error message", func(t *testing.T) {
		t.Parallel()

		err := errs.InvalidCapacityError{Capacity: 0, Minimum: 1, Policy: "active"}
		require.Equal(t, "capacity 0 must be greater than 1 for active caches", err.Error())
	})
}
