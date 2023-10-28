package executor

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.ATTEMPT.TRIGGER.EXECUTOR.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.ATTEMPT.TRIGGER.EXECUTOR.ALREAD_STARTED")
)
