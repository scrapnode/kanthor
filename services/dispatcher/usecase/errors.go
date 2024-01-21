package usecase

import "errors"

var (
	ErrNotConnected     = errors.New("DISPATCHER.USECASE.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("DISPATCHER.USECASE.ALREADY_CONNECTED.ERROR")
)
