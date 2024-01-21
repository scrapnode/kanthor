package cache

import "errors"

var (
	ErrAlreadyConnected = errors.New("INFRASTRUCTURE.CACHE.ALREADY_CONNECTED.ERROR")
	ErrNotConnected     = errors.New("INFRASTRUCTURE.CACHE.NOT_CONNECTED.ERROR")
	ErrEntryNotFound    = errors.New("INFRASTRUCTURE.CACHE.ENTRY.NOT_FOUND.ERROR")
)
