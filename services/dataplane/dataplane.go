package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dataplane/grpc/protos"
	"github.com/scrapnode/kanthor/usecases"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func New(conf *config.Config, logger logging.Logger, usecase usecases.Dataplane) services.Service {
	logger.With("service", "dataplane")
	return &dataplane{conf: conf, logger: logger, usecase: usecase}
}

type dataplane struct {
	conf    *config.Config
	logger  logging.Logger
	grpc    *grpc.Server
	usecase usecases.Dataplane
}

func (service *dataplane) Start(ctx context.Context) error {
	if err := service.usecase.Connect(ctx); err != nil {
		return err
	}

	service.grpc = grpc.NewServer()
	logger := service.logger.With("fn", "server")
	protos.RegisterDataplaneServer(service.grpc, &Server{usecase: service.usecase, logger: logger})
	reflection.Register(service.grpc)

	service.logger.Info("started")
	return nil
}

func (service *dataplane) Stop(ctx context.Context) error {
	service.grpc.GracefulStop()
	service.logger.Info("stopped")

	if err := service.usecase.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *dataplane) Run(ctx context.Context) error {
	addr := service.conf.Dataplane.Server.Addr

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	service.logger.Infow("running", "addr", addr)
	return service.grpc.Serve(listener)
}
