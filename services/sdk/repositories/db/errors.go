package db

import "errors"

var (
	ErrNotConnected     = errors.New("USECASES.SDK.REPOSITORIES.DB.CONNECTION.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("USECASES.SDK.REPOSITORIES.DB.CONNECTION.ALREADY_CONNECTED.ERROR")
)
