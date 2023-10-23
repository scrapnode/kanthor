package repos

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.SCHEDULER.REPOS.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.SCHEDULER.REPOS.CONNECTION.ALREADY_CONNECTED")
)
