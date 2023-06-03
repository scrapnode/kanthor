package msgbroker

import "errors"

var (
	ErrAlreadyConnected = errors.New("msgbroker is already connected")
	ErrNotConnected     = errors.New("msgbroker is not connected")
)
