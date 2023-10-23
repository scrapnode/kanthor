package scheduler

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.SCHEDULER.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.SCHEDULER.CONNECTION.ALREADY_CONNECTED")
)
