package consumer

import "errors"

var (
	ErrNotStarted     = errors.New("ATTEMPT.ENTRYPOINT.CONSUMER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("ATTEMPT.ENTRYPOINT.CONSUMER.ALREAD_STARTED")
)
