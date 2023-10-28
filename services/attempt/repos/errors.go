package repos

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.ATTEMPT.REPOS.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.ATTEMPT.REPOS.CONNECTION.ALREADY_CONNECTED")
)
