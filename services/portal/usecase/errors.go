package usecase

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.PORTAL.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.PORTAL.CONNECTION.ALREADY_CONNECTED")
)
