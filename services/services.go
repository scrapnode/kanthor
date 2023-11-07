package services

var (
	PORTAL                    = "portal"
	SDK                       = "sdk"
	SCHEDULER                 = "scheduler"
	DISPATCHER                = "dispatcher"
	STORAGE                   = "storage"
	ATTEMPT                   = "attempt"
	ATTEMPT_TRIGGER_PLANNER   = "attempt.trigger.planner"
	ATTEMPT_TRIGGER_EXECUTOR  = "attempt.trigger.executor"
	ATTEMPT_ENDEAVOR_PLANNER  = "attempt.endeavor.planner"
	ATTEMPT_ENDEAVOR_EXECUTOR = "attempt.endeavor.executor"
	SERVICES                  = []string{
		PORTAL,
		SDK,
		SCHEDULER,
		DISPATCHER,
		STORAGE,
		ATTEMPT,
		ATTEMPT_TRIGGER_PLANNER,
		ATTEMPT_TRIGGER_EXECUTOR,
		ATTEMPT_ENDEAVOR_PLANNER,
		ATTEMPT_ENDEAVOR_EXECUTOR,
	}
)
