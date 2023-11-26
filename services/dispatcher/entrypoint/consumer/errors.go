package consumer

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.DISPATCHER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.DISPATCHER.ALREAD_STARTED")
)
