package authorizator

import (
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) (Authorizator, error) {
	if conf.Engine == EngineNoop {
		return NewNoop(conf, logger), nil
	}

	if conf.Engine == EngineCasbin {
		return NewCasbin(conf, logger), nil
	}

	return nil, fmt.Errorf("authorizator: unknow engine")
}

type Authorizator interface {
	patterns.Connectable
	Enforce(sub, ws, obj, act string) (bool, error)
	SetupPermissions(role, tenant string, permissions [][]string) error
	GrantAccess(sub, role, tenant string) error
	Tenants(sub string) ([]string, error)
}
