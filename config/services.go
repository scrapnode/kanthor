package config

const (
	SERVICE_ALL        = "all"
	SERVICE_MIGRATION  = "migration"
	SERVICE_PORTAL_API = "portalapi"
	SERVICE_SDK_API    = "sdkapi"
	SERVICE_SCHEDULER  = "scheduler"
	SERVICE_DISPATCHER = "dispatcher"
	SERVICE_STORAGE    = "storage"
	SERVICE_ATTEMPT    = "attempt"
)

func IsValidServiceName(service string) bool {
	if service == SERVICE_ALL {
		return true
	}
	if service == SERVICE_MIGRATION {
		return true
	}
	if service == SERVICE_PORTAL_API {
		return true
	}
	if service == SERVICE_SDK_API {
		return true
	}
	if service == SERVICE_SCHEDULER {
		return true
	}
	if service == SERVICE_DISPATCHER {
		return true
	}
	if service == SERVICE_STORAGE {
		return true
	}
	if service == SERVICE_ATTEMPT {
		return true
	}

	return false
}
