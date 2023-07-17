package database

import (
	"errors"
)

var (
	ErrAlreadyConnected = errors.New("DATABASE.CONNECTION.ALREADY_CONNECTED")
	ErrNotConnected     = errors.New("DATABASE.CONNECTION.NOT_CONNECTED")
	ErrRecordNotFound   = errors.New("DATABASE.RECORD.NOT_FOUND")
)
