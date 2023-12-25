package db

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.PORTAL.repositories.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.PORTAL.repositories.CONNECTION.ALREADY_CONNECTED")
)
