package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
	grpccore "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	authenticator authenticator.Authenticator,
	meter metric.Meter,
	uc usecase.Controlplane,
) services.Service {
	logger = logger.With("gateway", "grpc")
	return &controlplane{conf: conf, logger: logger, authenticator: authenticator, meter: meter, uc: uc}
}

type controlplane struct {
	conf          *config.Config
	logger        logging.Logger
	gateway       *grpccore.Server
	authenticator authenticator.Authenticator
	meter         metric.Meter
	uc            usecase.Controlplane
}

func (service *controlplane) Start(ctx context.Context) error {
	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	service.gateway = grpc.New(service.logger, service.meter, service.authenticator)
	protos.RegisterWsServer(service.gateway, &ws{service: service})
	reflection.Register(service.gateway)

	service.logger.Info("started")
	return nil
}

func (service *controlplane) Stop(ctx context.Context) error {
	service.gateway.GracefulStop()
	service.logger.Info("stopped")

	if err := service.uc.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *controlplane) Run(ctx context.Context) error {
	addr := service.conf.Controlplane.GRPC.Addr

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	service.logger.Infow("running", "addr", addr)

	return service.gateway.Serve(listener)
}
