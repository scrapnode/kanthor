package scanner

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICES.RECOVERY.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICES.RECOVERY.ALREAD_STARTED")
)
