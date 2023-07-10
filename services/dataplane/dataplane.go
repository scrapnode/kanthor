package dataplane

import (
	"context"
	mwauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	mwrecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dataplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/dataplane"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"runtime/debug"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	uc usecase.Dataplane,
	authenticator authenticator.Authenticator,
	meter metric.Meter,
) services.Service {
	logger.With("service", "dataplane")
	return &dataplane{conf: conf, logger: logger, uc: uc, authenticator: authenticator, meter: meter}
}

type dataplane struct {
	conf          *config.Config
	logger        logging.Logger
	grpc          *grpc.Server
	uc            usecase.Dataplane
	authenticator authenticator.Authenticator
	meter         metric.Meter
}

func (service *dataplane) Start(ctx context.Context) error {
	if err := service.uc.Connect(ctx); err != nil {
		return err
	}

	service.grpc = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mwrecovery.UnaryServerInterceptor(mwrecovery.WithRecoveryHandler(func(p any) (err error) {
				service.logger.Errorw("recovered from panic", "panic", p, "stack", debug.Stack())
				return status.Error(codes.Internal, "oops, something went wrong")
			})),
			mwauth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {
				token, err := mwauth.AuthFromMD(ctx, "basic")
				if err != nil {
					return nil, err
				}

				account, err := service.authenticator.Verify(token)
				if err != nil {
					return nil, err
				}

				return context.WithValue(ctx, "kanthor.account", account), nil
			}),
		),
		grpc.ChainStreamInterceptor(
			mwrecovery.StreamServerInterceptor(mwrecovery.WithRecoveryHandler(func(p any) (err error) {
				service.logger.Errorw("recovered from panic", "panic", p, "stack", debug.Stack())
				return status.Error(codes.Internal, "oops, something went wrong")
			})),
			mwauth.StreamServerInterceptor(func(ctx context.Context) (context.Context, error) {
				token, err := mwauth.AuthFromMD(ctx, "basic")
				if err != nil {
					return nil, err
				}

				account, err := service.authenticator.Verify(token)
				if err != nil {
					return nil, err
				}

				return context.WithValue(ctx, "kanthor.account", account), nil
			}),
		))

	protos.RegisterMessageServer(service.grpc, &Message{service: service})
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
