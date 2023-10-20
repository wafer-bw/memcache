package memcache

import "errors"

var (
	ErrNilExpirerFunc  = errors.New("expirer function cannot be nil")
	ErrInvalidInterval = errors.New("interval must be greater than 0")
)
