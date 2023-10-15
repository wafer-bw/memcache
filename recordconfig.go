package memcache

import "time"

type RecordConfig struct {
	expireAt *time.Time
}

type RecordConfigOption func(*RecordConfig)

func WithTTL(d time.Duration) RecordConfigOption {
	return func(config *RecordConfig) {
		expireAt := time.Now().Add(d)
		config.expireAt = &expireAt
	}
}
