package memcache

type CacheConfig struct {
	passiveExpiration bool
}

type CacheConfigOption func(*CacheConfig) error

func WithPassiveExpiration() CacheConfigOption {
	return func(config *CacheConfig) error {
		config.passiveExpiration = true
		return nil
	}
}
