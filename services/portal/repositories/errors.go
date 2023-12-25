package repositories

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.ATTEMPT.REPOSITORIES.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.ATTEMPT.REPOSITORIES.CONNECTION.ALREADY_CONNECTED")
)
