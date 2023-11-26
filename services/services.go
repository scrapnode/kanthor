package services

var (
	ALL                       = "all"
	PORTAL                    = "portal"
	SDK                       = "sdk"
	SCHEDULER                 = "scheduler"
	DISPATCHER                = "dispatcher"
	STORAGE                   = "storage"
	ATTEMPT_TRIGGER_CLI       = "attempt.trigger.cli"
	ATTEMPT_TRIGGER_PLANNER   = "attempt.trigger.planner"
	ATTEMPT_TRIGGER_EXECUTOR  = "attempt.trigger.executor"
	ATTEMPT_ENDEAVOR_CLI      = "attempt.endeavor.cli"
	ATTEMPT_ENDEAVOR_PLANNER  = "attempt.endeavor.planner"
	ATTEMPT_ENDEAVOR_EXECUTOR = "attempt.endeavor.executor"
	SERVICES                  = append(
		[]string{
			PORTAL,
			SDK,
			SCHEDULER,
			DISPATCHER,
			STORAGE,
		},
		ATTEMPTS...,
	)
	ATTEMPTS = []string{
		ATTEMPT_TRIGGER_PLANNER,
		ATTEMPT_TRIGGER_EXECUTOR,
		ATTEMPT_ENDEAVOR_PLANNER,
		ATTEMPT_ENDEAVOR_EXECUTOR,
	}
)
