package datastore

import "errors"

var (
	ErrAlreadyConnected = errors.New("datastore: already connected")
	ErrNotConnected     = errors.New("datastore: not connected")
)
