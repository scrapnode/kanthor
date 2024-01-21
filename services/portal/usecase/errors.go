package usecase

import "errors"

var (
	ErrNotConnected     = errors.New("PORTAL.USECASE.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("PORTAL.USECASE.ALREADY_CONNECTED.ERROR")
)
