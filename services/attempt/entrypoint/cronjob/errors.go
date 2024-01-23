package cronjob

import "errors"

var (
	ErrNotStarted     = errors.New("RECOVERY.ENTRYPOINT.CRONJOB.NOT_STARTED")
	ErrAlreadyStarted = errors.New("RECOVERY.ENTRYPOINT.CRONJOB.ALREAD_STARTED")
)
