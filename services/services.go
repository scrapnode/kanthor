package services

var (
	PORTAL            = "portal"
	SDK               = "sdk"
	SCHEDULER         = "scheduler"
	DISPATCHER        = "dispatcher"
	STORAGE           = "storage"
	RECOVERY_CRONJOB  = "recovery.cronjob"
	RECOVERY_CONSUMER = "recovery.consumer"

	ALL      = "all"
	SERVICES = []string{
		PORTAL,
		SDK,
		SCHEDULER,
		DISPATCHER,
		STORAGE,
		RECOVERY_CRONJOB,
		RECOVERY_CONSUMER,
	}
)
