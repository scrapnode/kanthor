package repositories

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.SCHEDULER.repositories.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.SCHEDULER.repositories.CONNECTION.ALREADY_CONNECTED")
)
