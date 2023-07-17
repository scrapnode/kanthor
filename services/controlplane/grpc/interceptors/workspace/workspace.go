package workspace

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/stream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
	grpccore "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(logger logging.Logger, uc usecase.Controlplane) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		ctx, err = resolve(logger, uc, ctx)
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
		ctx, err := resolve(logger, uc, ss.Context())
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
	wsId := meta.Get("x-kanthor-workspace")
	// we will validate the workspace when we used it
	if wsId == "" {
		return ctx, nil
	}

	req := &usecase.WorkspaceGetReq{Id: wsId}
	res, err := uc.Workspace().Get(ctx, req)
	if err != nil {
		logger.Errorw(err.Error(), "workspace_id", wsId)
		return nil, status.Error(codes.Internal, "unable to get workspace")
	}

	logger.Debugw("resolve workspace", "workspace_id", res.Workspace.Id, "workspace_tier", res.Workspace.Tier)

	ctx = context.WithValue(ctx, authorizator.CtxAuthzWorkspace, res.Workspace)
	return ctx, nil
}
