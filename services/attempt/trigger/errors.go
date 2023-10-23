package trigger

import "errors"

var (
	ErrPlannerNotStarted      = errors.New("SERVICE.ATTEMPT.TRIGGER.PLANNER.NOT_STARTED")
	ErrPlannerAlreadyStarted  = errors.New("SERVICE.ATTEMPT.TRIGGER.PLANNER.ALREAD_STARTED")
	ErrExecutorNotStarted     = errors.New("SERVICE.ATTEMPT.TRIGGER.EXECUTOR.NOT_STARTED")
	ErrExecutorAlreadyStarted = errors.New("SERVICE.ATTEMPT.TRIGGER.EXECUTOR.ALREAD_STARTED")
)
