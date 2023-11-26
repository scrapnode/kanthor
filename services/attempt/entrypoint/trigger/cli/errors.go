package cli

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.ATTEMPT.TRIGGER.CLI.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.ATTEMPT.TRIGGER.CLI.ALREAD_STARTED")
)
