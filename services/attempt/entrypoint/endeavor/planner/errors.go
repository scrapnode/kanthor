package planner

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.ATTEMPT.ENDEAVOUR.PLANNER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.ATTEMPT.ENDEAVOUR.PLANNER.ALREAD_STARTED")
)
