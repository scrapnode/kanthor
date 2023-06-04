package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/dataplane/servers/grpc/protos"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func New(conf *config.Config, logger logging.Logger) (patterns.Runnable, error) {
	return &server{
		conf:   conf,
		logger: logger.With("component", "servers.grpc"),
	}, nil
}

type server struct {
	conf   *config.Config
	logger logging.Logger
	server *grpc.Server
}

func (s *server) Start(ctx context.Context) error {
	s.logger.Info("starting")

	s.server = grpc.NewServer()
	protos.RegisterMessageServer(s.server, &MessageServer{})
	reflection.Register(s.server)

	return nil
}

func (s *server) Stop(ctx context.Context) error {
	s.logger.Info("stopping")
	s.server.GracefulStop()
	return nil
}

func (s *server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.conf.Server.Addr)
	if err != nil {
		return err
	}

	s.logger.Infow("listening", "addr", s.conf.Server.Addr)
	return s.server.Serve(listener)
}
