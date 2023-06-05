package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/dataplane/servers/grpc/protos"
	"github.com/scrapnode/kanthor/dataplane/services"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func New(conf *config.Config, logger logging.Logger, services services.Services) patterns.Runnable {
	logger = logger.With("component", "servers.grpc")
	return &server{conf: conf, logger: logger, services: services}
}

type server struct {
	conf     *config.Config
	logger   logging.Logger
	services services.Services
	server   *grpc.Server
}

func (s *server) Start(ctx context.Context) error {
	if err := s.services.Connect(ctx); err != nil {
		return err
	}

	s.server = grpc.NewServer()
	protos.RegisterMessageServer(s.server, &MessageServer{service: s.services.Message()})
	reflection.Register(s.server)

	s.logger.Info("started")
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	s.server.GracefulStop()
	s.logger.Info("stopped")

	if err := s.services.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (s *server) Run(ctx context.Context) error {
	addr := s.conf.Dataplane.Server.Addr

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.logger.Infow("listening", "addr", addr)
	return s.server.Serve(listener)
}
