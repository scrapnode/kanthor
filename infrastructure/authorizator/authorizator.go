package authorizator

import (
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) Authorizator {
	return NewCasbin(conf, logger)
}

type Authorizator interface {
	patterns.Connectable
	Enforce(sub, ws, obj, act string) (bool, error)
	SetupPermissions(role, ws string, permissions [][]string) error
	GrantAccess(sub, role, ws string) error
	Workspaces(sub string) ([]string, error)
}
