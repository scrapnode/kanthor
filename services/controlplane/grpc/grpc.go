package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc"
	gatewayinterceptors "github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors/auth"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/interceptors"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/interceptors/authz"
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
	authorizator authorizator.Authorizator,
	meter metric.Meter,
	uc usecase.Controlplane,
) services.Service {
	logger = logger.With("gateway", "grpc")
	return &controlplane{
		conf:          conf,
		logger:        logger,
		authenticator: authenticator,
		meter:         meter,
		authorizator:  authorizator,
		uc:            uc,
	}
}

type controlplane struct {
	conf          *config.Config
	logger        logging.Logger
	gateway       *grpccore.Server
	authenticator authenticator.Authenticator
	authorizator  authorizator.Authorizator
	meter         metric.Meter
	uc            usecase.Controlplane
}

func (service *controlplane) Start(ctx context.Context) error {
	if err := service.authorizator.Connect(ctx); err != nil {
		return err
	}
	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	service.gateway = grpc.NewServer(
		gatewayinterceptors.WithRecovery(service.logger),
		gatewayinterceptors.WithMeasurement(service.meter),
		gatewayinterceptors.WithAuth(service.logger, service.authenticator, auth.DefaultPublic()),
		interceptors.WithWorkspace(service.logger, service.uc),
		interceptors.WithAuthz(service.logger, service.authorizator, authz.DefaultProtected()),
	)
	protos.RegisterAccountServer(service.gateway, &account{service: service})
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

	if err := service.authorizator.Disconnect(ctx); err != nil {
		service.logger.Error(err)
	}

	return nil
}

func (service *controlplane) Run(ctx context.Context) error {
	addr := service.conf.Controlplane.Gateway.GRPC.Addr

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	service.logger.Infow("running", "addr", addr)

	return service.gateway.Serve(listener)
}
