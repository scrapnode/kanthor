package authorizator

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/patterns"
)

var (
	HeaderAuthWorkspace = "x-authorization-workspace"
)

func New(conf *Config, logger logging.Logger) (Authorizator, error) {
	if conf.Engine == EngineCasbin {
		return NewCasbin(conf, logger)
	}

	return nil, fmt.Errorf("authorizator: unknown engine")
}

type Authorizator interface {
	patterns.Connectable
	Refresh(ctx context.Context) error

	Enforce(tenant, sub, obj, act string) (bool, error)
	Grant(tenant, sub, role string, permissions []Permission) error
	Tenants(sub string) ([]string, error)
	UsersOfTenant(tenant string) ([]string, error)
	UserPermissionsInTenant(tenant, sub string) ([]Permission, error)
}

type Permission struct {
	Role   string `json:"role,omitempty"`
	Object string `json:"object"`
	Action string `json:"action"`
}
