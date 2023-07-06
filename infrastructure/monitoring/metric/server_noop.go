package metric

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
)

func NewNoopServer(conf *Config, logger logging.Logger) services.Service {
	return &noopServer{conf: conf, logger: logger}
}

type noopServer struct {
	conf   *Config
	logger logging.Logger
}

func (server *noopServer) Start(ctx context.Context) error {
	server.logger.Info("connected")
	return nil
}

func (server *noopServer) Stop(ctx context.Context) error {
	server.logger.Info("disconnected")
	return nil
}

func (server *noopServer) Run(ctx context.Context) error {
	return nil
}
