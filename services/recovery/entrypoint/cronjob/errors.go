package cronjob

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICES.RECOVERY.CRONJOB.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICES.RECOVERY.CRONJOB.ALREAD_STARTED")
)
