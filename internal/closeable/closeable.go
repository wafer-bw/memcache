package closeable

type Closer struct {
	ch chan struct{}
}

func New() *Closer {
	return &Closer{ch: make(chan struct{})}
}

func (c *Closer) Ch() <-chan struct{} {
	return c.ch
}

func (c *Closer) Closed() bool {
	select {
	case <-c.ch:
		return true
	default:
		return false
	}
}

func (c *Closer) Close() {
	select {
	case <-c.ch:
		return
	default:
		close(c.ch)
	}
}
