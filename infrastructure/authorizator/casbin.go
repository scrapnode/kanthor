package authorizator

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	gocasbin "github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func NewCasbin(conf *Config, logger logging.Logger) Authorizator {
	logger = logger.With("authorizator", "casbin")
	return &casbin{conf: conf, logger: logger}
}

type casbin struct {
	conf   *Config
	logger logging.Logger

	mu      sync.Mutex
	adapter *gormadapter.Adapter
	client  *gocasbin.Enforcer
}

func (authorizator *casbin) Readiness() error {
	if authorizator.adapter == nil {
		return ErrNotConnected
	}

	var ok int
	tx := authorizator.adapter.GetDb().Raw("SELECT 1").Scan(&ok)
	if tx.Error != nil {
		return tx.Error
	}
	if ok != 1 {
		return ErrNotReady
	}

	return nil
}

func (authorizator *casbin) Liveness() error {
	if authorizator.adapter == nil {
		return ErrNotConnected
	}

	var ok int
	tx := authorizator.adapter.GetDb().Raw("SELECT 1").Scan(&ok)
	if tx.Error != nil {
		return tx.Error
	}
	if ok != 1 {
		return ErrNotLive
	}

	return nil
}

func (authorizator *casbin) Connect(ctx context.Context) error {
	authorizator.mu.Lock()
	defer authorizator.mu.Unlock()

	ns := authorizator.conf.Casbin.PolicyNamespace

	modelUrl, err := url.Parse(authorizator.conf.Casbin.ModelUri)
	if err != nil {
		return err
	}

	policyUrl, err := url.Parse(authorizator.conf.Casbin.PolicyUri)
	if err != nil {
		return err
	}
	databaseName := strings.ReplaceAll(policyUrl.Path, "/", "")
	tableName := fmt.Sprintf("kanthor_authz_%s", ns)

	adapter, err := gormadapter.NewAdapter(policyUrl.Scheme, authorizator.conf.Casbin.PolicyUri, databaseName, tableName, true)
	if err != nil {
		return err
	}
	authorizator.adapter = adapter

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

	authorizator.logger.Info("connected")
	return nil
}

func (authorizator *casbin) Disconnect(ctx context.Context) error {
	authorizator.mu.Lock()
	defer authorizator.mu.Unlock()

	if authorizator.adapter != nil {
		if err := authorizator.adapter.Close(); err != nil {
			authorizator.logger.Error(err)
		}
	}
	authorizator.adapter = nil

	authorizator.client = nil
	authorizator.logger.Info("disconnected")
	return nil
}

func (authorizator *casbin) Refresh(ctx context.Context) error {
	if err := authorizator.client.LoadModel(); err != nil {
		return err
	}
	if err := authorizator.client.LoadPolicy(); err != nil {
		return err
	}
	return nil
}

func (authorizator *casbin) Enforce(tenant, sub, obj, act string) (bool, error) {
	ok, explains, err := authorizator.client.EnforceEx(sub, tenant, obj, act)
	if err != nil {
		return false, err
	}

	if !ok {
		authorizator.logger.Warnw("permission denied", "sub", sub, "tenant", tenant, "obj", obj, "act", act, "explains", explains)
	}
	return ok, nil
}

func (authorizator *casbin) Grant(tenant, sub, role string, permissions []Permission) error {
	policies := [][]string{}
	for _, permission := range permissions {
		policies = append(policies, append([]string{role, tenant}, permission.Object, permission.Action))
	}
	// the returning boolean value indicates that whether we can add the entity or not
	// most time we could not add the new entity because it was exists already
	if _, err := authorizator.client.AddPolicies(policies); err != nil && !authorizator.IsUniqueViolation(err) {
		return err
	}

	// the returning boolean value indicates that whether we can add the entity or not
	// most time we could not add the new entity because it was exists already
	if _, err := authorizator.client.AddRoleForUserInDomain(sub, role, tenant); err != nil && !authorizator.IsUniqueViolation(err) {
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
		permissions = append(permissions, Permission{Role: policy[0], Object: policy[2], Action: policy[3]})
	}

	return permissions, nil
}

func (authorizator *casbin) IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgerrcode.UniqueViolation
	}

	return false
}
