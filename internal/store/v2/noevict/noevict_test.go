package noevict_test

import (
	"github.com/wafer-bw/memcache/internal/ports"
	"github.com/wafer-bw/memcache/internal/store/v2/noevict"
)

var _ ports.Storer[int, int] = (*noevict.Store[int, int])(nil)

// TODO: this.
