package sender

import (
	"errors"
)

var (
	// ErrUnsupportedTargetType is an error if user pass unsupported event target
	// type. Only supporting: reachable or addressable.
	ErrUnsupportedTargetType = errors.New("unsupported target type")
)
