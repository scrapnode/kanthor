package cronjob

import "errors"

var (
	ErrNotStarted     = errors.New("ATTEMPT.ENTRYPOINT.CRONJOB.NOT_STARTED")
	ErrAlreadyStarted = errors.New("ATTEMPT.ENTRYPOINT.CRONJOB.ALREAD_STARTED")
)
