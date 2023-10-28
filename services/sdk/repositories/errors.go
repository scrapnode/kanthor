package repositories

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.SDK.repositories.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.SDK.repositories.CONNECTION.ALREADY_CONNECTED")
)
