package compress

import (
	"errors"
)

var (
	ErrUnknown = errors.New("unknown compression format")
	ErrEmpty   = errors.New("no data to read")
)
