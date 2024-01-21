package database

import (
	"errors"
)

var (
	ErrNotReady         = errors.New("DATABASE.NOT_READY.ERROR")
	ErrNotLive          = errors.New("DATABASE.NOT_LIVE.ERROR")
	ErrAlreadyConnected = errors.New("DATABASE.ALREADY_CONNECTED.ERROR")
	ErrNotConnected     = errors.New("DATABASE.NOT_CONNECTED.ERROR")
	ErrRecordNotFound   = errors.New("DATABASE.RECORD.NOT_FOUND.ERROR")
)
