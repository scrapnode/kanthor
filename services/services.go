package services

var (
	PORTAL                   = "portal"
	SDK                      = "sdk"
	SCHEDULER                = "scheduler"
	DISPATCHER               = "dispatcher"
	STORAGE                  = "storage"
	ATTEMPT_TRIGGER_PLANNER  = "attempt.trigger.planner"
	ATTEMPT_TRIGGER_EXECUTOR = "attempt.trigger.executor"
	SERVICES                 = []string{
		PORTAL,
		SDK,
		SCHEDULER,
		DISPATCHER,
		STORAGE,
		ATTEMPT_TRIGGER_PLANNER,
		ATTEMPT_TRIGGER_EXECUTOR,
	}
)
