package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInterval = errors.New("provided interval must be greater than 0")
)

type InvalidCapacityError struct {
	Capacity int
	Minimum  int
	Policy   string
}

func (e InvalidCapacityError) Error() string {
	return fmt.Sprintf("capacity %d must be greater than %d for %s caches", e.Capacity, e.Minimum, e.Policy)
}
