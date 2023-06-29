package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dataplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/dataplane"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func New(conf *config.Config, logger logging.Logger, uc usecase.Dataplane) services.Service {
	logger.With("service", "dataplane")
	return &dataplane{conf: conf, logger: logger, uc: uc}
}

type dataplane struct {
	conf   *config.Config
	logger logging.Logger
	grpc   *grpc.Server
	uc     usecase.Dataplane
}

func (service *dataplane) Start(ctx context.Context) error {
	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	service.grpc = grpc.NewServer()
	logger := service.logger.With("fn", "server")
	protos.RegisterDataplaneServer(service.grpc, &Server{uc: service.uc, logger: logger})
	reflection.Register(service.grpc)

	service.logger.Info("started")
	return nil
}

func (service *dataplane) Stop(ctx context.Context) error {
	service.grpc.GracefulStop()
	service.logger.Info("stopped")

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *dataplane) Run(ctx context.Context) error {
	addr := service.conf.Dataplane.GRPC.Addr

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	service.logger.Infow("running", "addr", addr)
	return service.grpc.Serve(listener)
}
