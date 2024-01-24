package endeavor

import "errors"

var (
	ErrNotStarted     = errors.New("ATTEMPT.ENTRYPOINT.ENDEAVOR.NOT_STARTED")
	ErrAlreadyStarted = errors.New("ATTEMPT.ENTRYPOINT.ENDEAVOR.ALREAD_STARTED")
)
