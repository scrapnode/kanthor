package authorizator

import (
	"context"
	"fmt"
	gocasbin "github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	"net/url"
	"strings"
)

func NewCasbin(conf *Config, logger logging.Logger) Authorizator {
	logger = logger.With("authorizator", "casbin")
	return &casbin{conf: conf, logger: logger}
}

type casbin struct {
	conf   *Config
	logger logging.Logger

	watcher *watcher
	client  *gocasbin.Enforcer
}

func (authorizator *casbin) Connect(ctx context.Context) error {
	modelUrl, err := url.Parse(authorizator.conf.Casbin.ModelUri)
	if err != nil {
		return err
	}

	policyUrl, err := url.Parse(authorizator.conf.Casbin.PolicyUri)
	if err != nil {
		return err
	}
	databaseName := strings.ReplaceAll(policyUrl.Path, "/", "")
	tableName := fmt.Sprintf("kanthor_%s_rule", authorizator.conf.Casbin.PolicyNamespace)

	adapter, err := gormadapter.NewAdapter(policyUrl.Scheme, authorizator.conf.Casbin.PolicyUri, databaseName, tableName, true)
	if err != nil {
		return err
	}

	client, err := gocasbin.NewEnforcer(modelUrl.Host+modelUrl.Path, adapter)
	if err != nil {
		return err
	}
	if err := client.LoadModel(); err != nil {
		return err
	}
	if err := client.LoadPolicy(); err != nil {
		return err
	}
	authorizator.client = client

	authorizator.watcher = &watcher{
		nodeid:  utils.ID("casbin"),
		conf:    &authorizator.conf.Casbin.Watcher,
		logger:  authorizator.logger.With("casbin.watcher", "built-in"),
		subject: "kanthor.authorizator.casbin.watcher",
	}
	if err := authorizator.watcher.Connect(ctx); err != nil {
		return err
	}
	if err := authorizator.client.SetWatcher(authorizator.watcher); err != nil {
		return err
	}

	authorizator.logger.Info("connected")
	return nil
}

func (authorizator *casbin) Disconnect(ctx context.Context) error {
	if err := authorizator.watcher.Disconnect(ctx); err != nil {
		authorizator.logger.Error(err)
	}
	authorizator.watcher = nil

	authorizator.client = nil
	authorizator.logger.Info("disconnected")
	return nil
}

func (authorizator *casbin) Enforce(tenant, sub, obj, act string) (bool, error) {
	ok, explains, err := authorizator.client.EnforceEx(tenant, sub, obj, act)
	if err != nil {
		return false, err
	}

	if !ok {
		authorizator.logger.Warnw("permission denied", "sub", sub, "tenant", tenant, "obj", obj, "act", act, "explains", explains)
	}
	return ok, nil
}

func (authorizator *casbin) GrantPermissionsToRole(tenant, role string, permissions []Permission) error {
	policies := [][]string{}
	for _, permission := range permissions {
		policies = append(policies, append([]string{tenant, role}, permission.Object, permission.Action))
	}
	// the returning boolean value indicates that whether we can add the entity or not
	// most time we could not add the new entity because it was exists already
	_, err := authorizator.client.AddPolicies(policies)
	if err != nil {
		return err
	}

	return nil
}

func (authorizator *casbin) GrantRoleToSub(tenant, sub, role string) error {
	// the returning boolean value indicates that whether we can add the entity or not
	// most time we could not add the new entity because it was exists already
	_, err := authorizator.client.AddRoleForUserInDomain(tenant, sub, role)
	if err != nil {
		return err
	}

	return nil
}

func (authorizator *casbin) Tenants(sub string) ([]string, error) {
	tenants, err := authorizator.client.GetDomainsForUser(sub)
	if err != nil {
		return nil, err
	}

	if tenants == nil {
		tenants = []string{}
	}
	return tenants, nil
}

func (authorizator *casbin) UsersOfTenant(tenant string) ([]string, error) {
	users := authorizator.client.GetAllUsersByDomain(tenant)

	if users == nil {
		users = []string{}
	}
	return users, nil
}

func (authorizator *casbin) UserPermissionsInTenant(tenant, sub string) ([]Permission, error) {
	permissions := []Permission{}

	policies := authorizator.client.GetPermissionsForUserInDomain(sub, tenant)
	for _, policy := range policies {
		permissions = append(permissions, Permission{Object: policy[0], Action: policy[1]})
	}

	return permissions, nil
}
