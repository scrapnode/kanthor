package db

import "errors"

var (
	ErrNotConnected     = errors.New("PORTAL.REPOSITORIES.DB.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected = errors.New("PORTAL.REPOSITORIES.DB.ALREADY_CONNECTED.ERROR")
)
