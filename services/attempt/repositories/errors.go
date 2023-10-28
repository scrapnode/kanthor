package repositories

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.ATTEMPT.repositories.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.ATTEMPT.repositories.CONNECTION.ALREADY_CONNECTED")
)
