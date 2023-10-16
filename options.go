package memcache

type Option[K comparable, V any] func(*Cache[K, V]) error

func WithPassiveExpiration[K comparable, V any]() Option[K, V] {
	return func(c *Cache[K, V]) error {
		c.passiveExpiration = true
		return nil
	}
}

// TODO: WithExpirer

// TODO: WithEvictor
