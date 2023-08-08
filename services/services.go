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
