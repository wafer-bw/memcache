package allkeyslru_test

import (
	"github.com/wafer-bw/memcache/internal/ports"
	"github.com/wafer-bw/memcache/internal/store/v2/allkeyslru"
)

var _ ports.Storer[int, int] = (*allkeyslru.Store[int, int])(nil)

// TODO: this.
