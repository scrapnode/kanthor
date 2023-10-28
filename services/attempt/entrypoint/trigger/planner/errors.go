package trigger

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.ATTEMPT.TRIGGER.PLANNER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.ATTEMPT.TRIGGER.PLANNER.ALREAD_STARTED")
)
