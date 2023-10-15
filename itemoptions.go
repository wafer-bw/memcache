package memcache

import "time"

type ItemConfig struct {
	expireAt *time.Time
}

type ItemConfigOption func(*ItemConfig)

func WithTTL(d time.Duration) ItemConfigOption {
	return func(itemConfig *ItemConfig) {
		expireAt := time.Now().Add(d)
		itemConfig.expireAt = &expireAt
	}
}
