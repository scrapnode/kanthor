package entrypoint

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICES.STORAGE.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICES.STORAGE.ALREAD_STARTED")
)
