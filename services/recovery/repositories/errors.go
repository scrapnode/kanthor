package repositories

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.ATTEMPT.REPOSITORIES.CONNECTION.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("USECASES.ATTEMPT.REPOSITORIES.CONNECTION.ALREADY_CONNECTED.ERROR")
)
