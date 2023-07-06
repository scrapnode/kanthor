package metric

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"net"
	"net/http"
)

func NewHttpServer(conf *Config, logger logging.Logger, handler http.Handler) services.Service {
	return &httpServer{conf: conf, logger: logger, handler: handler}
}

type httpServer struct {
	conf   *Config
	logger logging.Logger

	handler http.Handler
	server  *http.Server
}

func (server *httpServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/healthz", &healthz{logger: server.logger})
	mux.Handle("/metrics", server.handler)
	server.server = &http.Server{Handler: mux}

	server.logger.Info("connected")
	return nil
}

func (server *httpServer) Stop(ctx context.Context) error {
	if err := server.server.Shutdown(ctx); err != nil {
		return err
	}

	server.logger.Info("disconnected")
	return nil
}

func (server *httpServer) Run(ctx context.Context) error {
	addr := server.conf.Exporter.Addr

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server.logger.Infow("running", "addr", addr)
	err = server.server.Serve(listener)

	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

type healthz struct {
	logger logging.Logger
}

func (handler healthz) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("healthy"))
	if err != nil {
		handler.logger.Errorw("unable to response client", "error", err.Error(), "handler", "healthz")
	}
}
