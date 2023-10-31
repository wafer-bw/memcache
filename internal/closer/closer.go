package closer

import "sync"

type Closer struct {
	mu sync.RWMutex
	ch chan struct{}
}

func New() *Closer {
	return &Closer{ch: make(chan struct{})}
}

func (c *Closer) WaitClosed() <-chan struct{} {
	return c.ch
}

func (c *Closer) Closed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	select {
	case <-c.ch:
		return true
	default:
		return false
	}
}

func (c *Closer) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.ch:
		return
	default:
		close(c.ch)
	}
}
