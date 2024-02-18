package closer

import "sync"

type Closer struct {
	mu     sync.RWMutex // TODO: is this mutex necessary?
	once   sync.Once
	closed bool
}

func New() *Closer {
	return &Closer{
		mu: sync.RWMutex{},
	}
}

func (c *Closer) Close() {
	c.once.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.closed = true
	})
}

func (c *Closer) Closed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}
