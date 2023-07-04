package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"net"
	"net/http"
)

func New(conf *config.Server, logger logging.Logger) services.Service {
	logger.With("service", "metrics")
	return &metrics{conf: conf}
}

type metrics struct {
	conf   *config.Server
	logger logging.Logger

	server *http.Server
}

func (service *metrics) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	service.server = &http.Server{Handler: mux}

	service.logger.Info("connected")
	return nil
}

func (service *metrics) Stop(ctx context.Context) error {
	if err := service.server.Shutdown(ctx); err != nil {
		return err
	}

	service.logger.Info("disconnected")
	return nil
}

func (service *metrics) Run(ctx context.Context) error {
	addr := service.conf.Addr

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	service.logger.Infow("running", "addr", addr)
	return service.server.Serve(listener)
}
