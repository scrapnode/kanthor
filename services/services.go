package services

import "github.com/scrapnode/kanthor/infrastructure/patterns"

type Service interface {
	patterns.Runnable
}

const (
	ALL        = "all"
	PORTAL     = "portal"
	SDK        = "sdk"
	SCHEDULER  = "scheduler"
	DISPATCHER = "dispatcher"
	MIGRATION  = "migration"
)
