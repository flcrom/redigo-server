package store

import (
	"errors"
)

var ErrWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
var ErrValNotInt = errors.New("ERR value is not an integer or out of range")
