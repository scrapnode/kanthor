package usecase

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.ATTEMPT.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.ATTEMPT.CONNECTION.ALREADY_CONNECTED")
)
