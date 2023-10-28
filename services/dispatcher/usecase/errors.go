package usecase

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.DISPATCHER.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.DISPATCHER.CONNECTION.ALREADY_CONNECTED")
)
