package repositories

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.STORAGE.repositories.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.STORAGE.repositories.CONNECTION.ALREADY_CONNECTED")
)
