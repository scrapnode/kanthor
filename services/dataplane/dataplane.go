package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dataplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/dataplane"
	grpccore "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	uc usecase.Dataplane,
	authenticator authenticator.Authenticator,
	meter metric.Meter,
) services.Service {
	logger = logger.With("service", "dataplane")
	return &dataplane{conf: conf, logger: logger, uc: uc, authenticator: authenticator, meter: meter}
}

type dataplane struct {
	conf          *config.Config
	logger        logging.Logger
	gateway       *grpccore.Server
	uc            usecase.Dataplane
	authenticator authenticator.Authenticator
	meter         metric.Meter
}

func (service *dataplane) Start(ctx context.Context) error {
	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	service.gateway = grpc.New(service.logger, service.authenticator)
	protos.RegisterMessageServer(service.gateway, &Message{service: service})
	reflection.Register(service.gateway)

	service.logger.Info("started")
	return nil
}

func (service *dataplane) Stop(ctx context.Context) error {
	service.gateway.GracefulStop()
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

	return service.gateway.Serve(listener)
}
