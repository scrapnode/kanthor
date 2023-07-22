package application

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/stream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	usecase "github.com/scrapnode/kanthor/usecases/dataplane"
	grpccore "google.golang.org/grpc"
)

func UnaryServerInterceptor(logger logging.Logger, uc usecase.Dataplane) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		ctx, err = resolve(logger.With("method", info.FullMethod), uc, ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(logger logging.Logger, uc usecase.Dataplane) grpccore.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpccore.ServerStream,
		info *grpccore.StreamServerInfo,
		handler grpccore.StreamHandler,
	) error {
		ctx, err := resolve(logger.With("method", info.FullMethod), uc, ss.Context())
		if err != nil {
			return err
		}

		wrapped := stream.WrapServerStream(ss)
		wrapped.WrappedContext = ctx
		return handler(srv, wrapped)
	}
}

func resolve(logger logging.Logger, uc usecase.Dataplane, ctx context.Context) (context.Context, error) {
	meta := gateway.ExtractIncoming(ctx)
	appId := meta.Get("x-kanthor-app")
	// we don't need to ensure the existing of workspace at this place
	// because we have to check the value of ctx.Value(usecase.CtxWorkspace) later whenever we use it
	if appId == "" {
		return ctx, nil
	}

	req := &usecase.ApplicationGetReq{Id: appId}
	res, err := uc.Application().Get(ctx, req)
	if err != nil {
		logger.Errorw(err.Error(), "app_id", appId)
		return ctx, gateway.Err404("APPLICATION.NOT_FOUND")
	}

	logger.Debugw("resolve application", "app_id", appId, "workspace_id", res.Workspace.Id, "workspace_tier", res.Workspace.Tier)

	ctx = context.WithValue(ctx, usecase.CtxWorkspace, res.Workspace)
	ctx = context.WithValue(ctx, usecase.CtxApplication, res.Application)
	return ctx, nil
}
