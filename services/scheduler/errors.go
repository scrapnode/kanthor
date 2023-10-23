package scheduler

import "errors"

var (
	ErrNotStarted     = errors.New("SERVICE.SCHEDULER.NOT_STARTED")
	ErrAlreadyStarted = errors.New("SERVICE.SCHEDULER.ALREAD_STARTED")
)
