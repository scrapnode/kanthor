package db

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.PORTAL.repositories.CONNECTION.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("USECASES.PORTAL.repositories.CONNECTION.ALREADY_CONNECTED.ERROR")
)
