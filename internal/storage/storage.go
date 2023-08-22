package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrTaskExists  = errors.New("url exists")
)
