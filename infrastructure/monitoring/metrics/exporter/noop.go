package exporter

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"net/http"
)

func NewNoop(conf *Config, logger logging.Logger, handler http.Handler) Exporter {
	return &noop{conf: conf, logger: logger}
}

type noop struct {
	conf   *Config
	logger logging.Logger
}

func (exporter *noop) Start(ctx context.Context) error {
	exporter.logger.Info("started")
	return nil
}

func (exporter *noop) Stop(ctx context.Context) error {
	exporter.logger.Info("stopped")
	return nil
}

func (exporter *noop) Run(ctx context.Context) error {
	exporter.logger.Infow("running", "addr", exporter.conf.Addr)
	return nil
}
