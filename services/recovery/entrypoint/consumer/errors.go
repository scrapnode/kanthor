package consumer

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICES.RECOVERY.SCANNER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICES.RECOVERY.SCANNER.ALREAD_STARTED")
)
