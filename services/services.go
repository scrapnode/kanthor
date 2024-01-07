package services

var (
	PORTAL           = "portal"
	SDK              = "sdk"
	SCHEDULER        = "scheduler"
	DISPATCHER       = "dispatcher"
	STORAGE          = "storage"
	RECOVERY_SCANNER = "recovery.scanner"

	ALL      = "all"
	SERVICES = []string{
		PORTAL,
		SDK,
		SCHEDULER,
		DISPATCHER,
		STORAGE,
		RECOVERY_SCANNER,
	}
)
