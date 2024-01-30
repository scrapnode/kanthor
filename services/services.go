package services

var (
	PORTAL            = "portal"
	SDK               = "sdk"
	SCHEDULER         = "scheduler"
	DISPATCHER        = "dispatcher"
	STORAGE           = "storage"
	RECOVERY_CRONJOB  = "recovery.cronjob"
	RECOVERY_CONSUMER = "recovery.consumer"
	ATTEMPT_CRONJOB   = "attempt.cronjob"
	ATTEMPT_CONSUMER  = "attempt.consumer"
	ATTEMPT_TRIGGER   = "attempt.trigger"
	ATTEMPT_SELECTOR  = "attempt.selector"
	ATTEMPT_ENDEAVOR  = "attempt.endeavor"

	ALL      = "all"
	SERVICES = []string{
		PORTAL,
		SDK,
		SCHEDULER,
		DISPATCHER,
		STORAGE,
		RECOVERY_CRONJOB,
		RECOVERY_CONSUMER,
		ATTEMPT_CRONJOB,
		ATTEMPT_CONSUMER,
		ATTEMPT_TRIGGER,
		ATTEMPT_SELECTOR,
		ATTEMPT_ENDEAVOR,
	}
	SERVICE_RECOVERY = []string{
		RECOVERY_CRONJOB,
		RECOVERY_CONSUMER,
	}
	SERVICE_ATTEMPT = []string{
		ATTEMPT_CRONJOB,
		ATTEMPT_CONSUMER,
		ATTEMPT_TRIGGER,
		ATTEMPT_SELECTOR,
		ATTEMPT_ENDEAVOR,
	}
)
