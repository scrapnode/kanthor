package usecase

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.STORAGE.CONNECTION.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("USECASES.STORAGE.CONNECTION.ALREADY_CONNECTED.ERROR")
)
