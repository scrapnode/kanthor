package datastore

import "errors"

var (
	ErrNotReady         = errors.New("DATASTORE.NOT_READY")
	ErrNotLive          = errors.New("DATASTORE.NOT_LIVE")
	ErrAlreadyConnected = errors.New("DATASTORE.CONNECTION.ALREADY_CONNECTED")
	ErrNotConnected     = errors.New("DATASTORE.CONNECTION.NOT_CONNECTED")
	ErrRecordNotFound   = errors.New("DATASTORE.RECORD.NOT_FOUND")
)
