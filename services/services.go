package services

import "github.com/scrapnode/kanthor/infrastructure/patterns"

type Service interface {
	patterns.Runnable
}

const (
	DATAPLANE  = "dataplane"
	SCHEDULER  = "scheduler"
	DISPATCHER = "dispatcher"
	MIGRATION  = "migration"
)
