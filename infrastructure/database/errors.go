package database

import "errors"

var (
	ErrAlreadyConnected = errors.New("database: already connected")
	ErrNotConnected     = errors.New("database: not connected")
)
