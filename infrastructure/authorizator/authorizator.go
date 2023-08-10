package authorizator

import (
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

var (
	HeaderWorkspace = "x-kanthor-workspace-id"
)

func New(conf *Config, logger logging.Logger) (Authorizator, error) {
	if conf.Engine == EngineCasbin {
		return NewCasbin(conf, logger), nil
	}

	return nil, fmt.Errorf("authorizator: unknown engine")
}

type Authorizator interface {
	patterns.Connectable
	Enforce(sub, tenant, obj, act string) (bool, error)
	GrantPermissionsToRole(tenant, role string, permissions []Permission) error
	GrantRoleToSub(tenant, role, sub string) error
	Tenants(sub string) ([]string, error)
	UsersOfTenant(tenant string) ([]string, error)
	UserPermissionsInTenant(tenant, sub string) ([]Permission, error)
}

type Permission struct {
	Object string `json:"object"`
	Action string `json:"action"`
}
