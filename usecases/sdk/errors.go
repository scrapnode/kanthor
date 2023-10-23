package sdk

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.SDK.CONNECTION.NOT_CONNECTED")
	ErrAlreadyConnected = errors.New("USECASES.SDK.CONNECTION.ALREADY_CONNECTED")
)
