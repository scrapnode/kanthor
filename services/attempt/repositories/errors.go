package repositories

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.PORTAL.REPOSITORIES.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.PORTAL.REPOSITORIES.CONNECTION.ALREADY_CONNECTED")
)
