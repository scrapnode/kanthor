package authorizator

import (
	"context"
	gocasbin "github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	"net/url"
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
	adapter, err := gormadapter.NewAdapter(policyUrl.Scheme, authorizator.conf.Casbin.PolicyUri, true)
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

func (authorizator *casbin) Enforce(sub, tenant, obj, act string) (bool, error) {
	ok, explains, err := authorizator.client.EnforceEx(sub, tenant, obj, act)
	if err != nil {
		return false, err
	}

	if !ok {
		authorizator.logger.Warnw("permission denied", "sub", sub, "tenant", tenant, "obj", obj, "act", act, "explains", explains)
	}
	return ok, nil
}

func (authorizator *casbin) SetupPermissions(role, tenant string, permissions [][]string) error {
	var policies [][]string
	for _, permission := range permissions {
		policies = append(policies, append([]string{role, tenant}, permission...))
	}
	// the returning boolean value indicates that whether we can add the entity or not
	// most time we could not add the new entity because it was exists already
	_, err := authorizator.client.AddPolicies(policies)
	if err != nil {
		return err
	}

	return nil
}

func (authorizator *casbin) GrantAccess(sub, role, tenant string) error {
	// the returning boolean value indicates that whether we can add the entity or not
	// most time we could not add the new entity because it was exists already
	_, err := authorizator.client.AddRoleForUserInDomain(sub, role, tenant)
	if err != nil {
		return err
	}

	return nil
}

func (authorizator *casbin) Tenants(sub string) ([]string, error) {
	return authorizator.client.GetDomainsForUser(sub)
}
