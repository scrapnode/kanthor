package usecase

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.SCHEDULER.CONNECTION.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("USECASES.SCHEDULER.CONNECTION.ALREADY_CONNECTED.ERROR")
)
