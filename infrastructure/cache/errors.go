package cache

import "errors"

var (
	ErrAlreadyConnected = errors.New("cache: already connected")
	ErrNotConnected     = errors.New("cache: not connected")
	ErrEntryNotFound    = errors.New("cache: entry is not found")
)
