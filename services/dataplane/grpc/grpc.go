package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc"
	gatewayinterceptors "github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/auth"
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
	authenticator authenticator.Authenticator,
	meter metric.Meter,
	uc usecase.Dataplane,
) services.Service {
	logger = logger.With("gateway", "grpc")
	return &dataplane{conf: conf, logger: logger, authenticator: authenticator, meter: meter, uc: uc}
}

type dataplane struct {
	conf          *config.Config
	logger        logging.Logger
	gateway       *grpccore.Server
	authenticator authenticator.Authenticator
	meter         metric.Meter
	uc            usecase.Dataplane
}

func (service *dataplane) Start(ctx context.Context) error {
	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	service.gateway = grpc.NewServer(
		gatewayinterceptors.WithRecovery(service.logger),
		gatewayinterceptors.WithMeasurement(service.meter),
		gatewayinterceptors.WithAuth(service.logger, service.authenticator, auth.DefaultPublic()),
	)
	protos.RegisterMsgServer(service.gateway, &msg{service: service})
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
	addr := service.conf.Dataplane.Gateway.GRPC.Addr

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	service.logger.Infow("running", "addr", addr)

	return service.gateway.Serve(listener)
}
