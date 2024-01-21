package datastore

import "errors"

var (
	ErrNotReady         = errors.New("DATASTORE.NOT_READY.ERROR")
	ErrNotLive          = errors.New("DATASTORE.NOT_LIVE.ERROR")
	ErrAlreadyConnected = errors.New("DATASTORE.ALREADY_CONNECTED.ERROR")
	ErrNotConnected     = errors.New("DATASTORE.NOT_CONNECTED.ERROR")
	ErrRecordNotFound   = errors.New("DATASTORE.RECORD.NOT_FOUND.ERROR")
)
