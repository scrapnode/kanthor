package database

import (
	"errors"
)

var (
	ErrNotReady         = errors.New("DATABASE.NOT_READY")
	ErrNotLive          = errors.New("DATABASE.NOT_LIVE")
	ErrAlreadyConnected = errors.New("DATABASE.CONNECTION.ALREADY_CONNECTED")
	ErrNotConnected     = errors.New("DATABASE.CONNECTION.NOT_CONNECTED")
	ErrRecordNotFound   = errors.New("DATABASE.RECORD.NOT_FOUND")
)
