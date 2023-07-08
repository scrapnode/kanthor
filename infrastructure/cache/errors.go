package cache

import "errors"

var (
	ErrAlreadyConnected = errors.New("CACHE.CONNECTION.ALREADY_CONNECTED")
	ErrNotConnected     = errors.New("CACHE.CONNECTION.NOT_CONNECTED")
	ErrEntryNotFound    = errors.New("CACHE.ENTRY.NOT_FOUND")
)
