package ds

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.STORAGE.repositories.CONNECTION.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("USECASES.STORAGE.repositories.CONNECTION.ALREADY_CONNECTED.ERROR")
)
