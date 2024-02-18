package closeable

import "sync"

type Close struct {
	mu     sync.RWMutex // TODO: is this mutex necessary?
	once   sync.Once
	closed bool
}

func New() *Close {
	return &Close{
		mu: sync.RWMutex{},
	}
}

func (c *Close) Close() {
	c.once.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.closed = true
	})
}

func (c *Close) Closed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}
