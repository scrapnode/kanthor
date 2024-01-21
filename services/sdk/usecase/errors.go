package usecase

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.SDK.CONNECTION.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("USECASES.SDK.CONNECTION.ALREADY_CONNECTED.ERROR")
)
