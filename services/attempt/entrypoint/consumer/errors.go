package consumer

import "errors"

var (
	ErrNotStarted     = errors.New("RECOVERY.ENTRYPOINT.CONSUMER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("RECOVERY.ENTRYPOINT.CONSUMER.ALREAD_STARTED")
)
