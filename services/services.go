package services

import "github.com/scrapnode/kanthor/infrastructure/patterns"

type Service interface {
	patterns.Runnable
}

const (
	ALL        = "all"
	MIGRATION  = "migration"
	PORTAL_API = "portalapi"
	SDK_API    = "sdkapi"
	SCHEDULER  = "scheduler"
	DISPATCHER = "dispatcher"
	STORAGE    = "storage"
)

func Valid(service string) bool {
	if service == ALL {
		return true
	}
	if service == MIGRATION {
		return true
	}
	if service == PORTAL_API {
		return true
	}
	if service == SDK_API {
		return true
	}
	if service == SCHEDULER {
		return true
	}
	if service == DISPATCHER {
		return true
	}
	if service == STORAGE {
		return true
	}

	return false
}
