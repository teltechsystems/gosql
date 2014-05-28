package gosql

import (
	"errors"
)

var (
	MissingDatabase = errors.New("Missing database association")
)
