package storage

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.STORAGE.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.STORAGE.CONNECTION.ALREADY_CONNECTED")
)
