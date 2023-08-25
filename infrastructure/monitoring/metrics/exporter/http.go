package exporter

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"net/http"
)

func NewHttp(conf *Config, logger logging.Logger, handler http.Handler) Exporter {
	return &httpe{conf: conf, logger: logger, handler: handler}
}

type httpe struct {
	conf   *Config
	logger logging.Logger

	handler http.Handler
	server  *http.Server
}

func (exporter *httpe) Start(ctx context.Context) error {
	exporter.server = &http.Server{
		Addr:    exporter.conf.Addr,
		Handler: exporter.handler,
	}

	exporter.logger.Info("started")
	return nil
}

func (exporter *httpe) Stop(ctx context.Context) error {
	exporter.logger.Info("stopped")

	if err := exporter.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (exporter *httpe) Run(ctx context.Context) error {
	exporter.logger.Infow("running", "addr", exporter.conf.Addr)

	if err := exporter.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		exporter.logger.Error(err)
	}

	return nil
}
