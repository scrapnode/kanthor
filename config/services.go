package config

import "slices"

var (
	SERVICE_ALL                      = "all"
	SERVICE_PORTAL_API               = "portalapi"
	SERVICE_SDK_API                  = "sdkapi"
	SERVICE_SCHEDULER                = "scheduler"
	SERVICE_DISPATCHER               = "dispatcher"
	SERVICE_STORAGE                  = "storage"
	SERVICE_ATTEMPT_TRIGGER_PLANNER  = "attempt.trigger.planner"
	SERVICE_ATTEMPT_TRIGGER_EXECUTOR = "attempt.trigger.executor"
	SERVICES                         = []string{
		SERVICE_PORTAL_API,
		SERVICE_SDK_API,
		SERVICE_SCHEDULER,
		SERVICE_DISPATCHER,
		SERVICE_STORAGE,
		SERVICE_ATTEMPT_TRIGGER_PLANNER,
		SERVICE_ATTEMPT_TRIGGER_EXECUTOR,
	}
)

func IsValidServiceName(service string) bool {
	if service == SERVICE_ALL {
		return true
	}
	return slices.Contains(SERVICES, service)
}
