package trigger

import "errors"

var (
	ErrNotStarted     = errors.New("ATTEMPT.ENTRYPOINT.TRIGGER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("ATTEMPT.ENTRYPOINT.TRIGGER.ALREAD_STARTED")
)
