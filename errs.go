package memcache

import "errors"

var (
	ErrInvalidInterval = errors.New("interval must be greater than 0")
	ErrInvalidCapacity = errors.New("capacity must be greater than 1")
)
