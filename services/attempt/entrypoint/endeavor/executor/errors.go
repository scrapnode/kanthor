package executor

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.ATTEMPT.ENDEAVOUR.EXECUTOR.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.ATTEMPT.ENDEAVOUR.EXECUTOR.ALREAD_STARTED")
)
