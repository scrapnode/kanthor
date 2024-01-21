package consumer

import "errors"

var (
	ErrNotStarted     = errors.New("DISPATCHER.ENTRYPOINT.CONSUMER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("DISPATCHER.ENTRYPOINT.CONSUMER.ALREAD_STARTED")
)
