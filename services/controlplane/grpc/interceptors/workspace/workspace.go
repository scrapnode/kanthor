package workspace

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/stream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
	grpccore "google.golang.org/grpc"
)

func UnaryServerInterceptor(logger logging.Logger, uc usecase.Controlplane) grpccore.UnaryServerInterceptor {
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

func StreamServerInterceptor(logger logging.Logger, uc usecase.Controlplane) grpccore.StreamServerInterceptor {
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

func resolve(logger logging.Logger, uc usecase.Controlplane, ctx context.Context) (context.Context, error) {
	meta := gateway.ExtractIncoming(ctx)
	wsId := meta.Get("x-kanthor-ws")
	// we don't need to ensure the existing of workspace at this place
	// because we have to check the value of ctx.Value(usecase.CtxWorkspace) later whenever we use it
	if wsId == "" {
		return ctx, nil
	}

	req := &usecase.WorkspaceGetReq{Id: wsId}
	res, err := uc.Workspace().Get(ctx, req)
	if err != nil {
		logger.Errorw(err.Error(), "workspace_id", wsId)
		return ctx, gateway.Err404("WORKSPACE.NOT_FOUND")
	}

	logger.Debugw("resolve workspace", "workspace_id", res.Workspace.Id, "workspace_tier", res.Workspace.Tier)

	ctx = context.WithValue(ctx, usecase.CtxWorkspace, res.Workspace)
	return ctx, nil
}
