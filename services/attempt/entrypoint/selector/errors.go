package selector

import "errors"

var (
	ErrNotStarted     = errors.New("ATTEMPT.ENTRYPOINT.SELECTOR.NOT_STARTED")
	ErrAlreadyStarted = errors.New("ATTEMPT.ENTRYPOINT.SELECTOR.ALREAD_STARTED")
)
