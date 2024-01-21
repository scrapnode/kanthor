package usecase

import "errors"

var (
	ErrNotConnected     = errors.New("RECOVERY.USECASE.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("RECOVERY.USECASE.ALREADY_CONNECTED.ERROR")
)
