package authorizator

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

func NewNoop(conf *Config, logger logging.Logger) Authorizator {
	logger = logger.With("authorizator", "noop")
	return &noop{conf: conf, logger: logger}
}

type noop struct {
	conf   *Config
	logger logging.Logger
}

func (authorizator *noop) Connect(ctx context.Context) error {
	authorizator.logger.Info("connected")
	return nil
}

func (authorizator *noop) Disconnect(ctx context.Context) error {
	authorizator.logger.Info("disconnected")
	return nil
}

func (authorizator *noop) Enforce(sub, ws, obj, act string) (bool, error) {
	return true, nil
}

func (authorizator *noop) SetupPermissions(role, tenant string, permissions [][]string) error {
	return nil
}

func (authorizator *noop) GrantAccess(sub, role, tenant string) error {
	return nil
}

func (authorizator *noop) Tenants(sub string) ([]string, error) {
	return []string{}, nil
}
