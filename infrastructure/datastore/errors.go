package datastore

import "errors"

var (
	ErrAlreadyConnected = errors.New("DATASTORE.CONNECTION.ALREADY_CONNECTED")
	ErrNotConnected     = errors.New("DATASTORE.CONNECTION.NOT_CONNECTED")
)
