package memcache

import "time"

type ItemConfig struct {
	expireAt *time.Time
}

type ItemOption func(*ItemConfig)

func WithTTL(d time.Duration) ItemOption {
	return func(config *ItemConfig) {
		expireAt := time.Now().Add(d)
		config.expireAt = &expireAt
	}
}
