package repos

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.PORTAL.REPOS.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.PORTAL.REPOS.CONNECTION.ALREADY_CONNECTED")
)
