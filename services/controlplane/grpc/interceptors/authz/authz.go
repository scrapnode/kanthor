package authz

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/enforcer"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/metadata"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/stream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
	grpccore "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func UnaryServerInterceptor(
	logger logging.Logger,
	uc usecase.Controlplane,
	engine enforcer.Enforcer,
) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		ctx, err = authorize(logger, uc, engine, ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(
	logger logging.Logger,
	uc usecase.Controlplane,
	engine enforcer.Enforcer,
) grpccore.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpccore.ServerStream,
		info *grpccore.StreamServerInfo,
		handler grpccore.StreamHandler,
	) error {
		ctx, err := authorize(logger, uc, engine, ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		wrapped := stream.WrapServerStream(ss)
		wrapped.WrappedContext = ctx
		return handler(srv, wrapped)
	}
}

func authorize(
	logger logging.Logger,
	uc usecase.Controlplane,
	engine enforcer.Enforcer,
	ctx context.Context,
	method string,
) (context.Context, error) {
	acc := authenticator.AccountFromContext(ctx)
	res, err := ws(uc, ctx, acc)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	obj, act, err := action(method)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	ok, err := engine.Enforce(res.Privilege.AccountRole, obj, act)
	if err != nil {
		logger.Errorw(err.Error(), "workspace_id", res.Workspace.Id, "account_sub", acc.Sub)
		return nil, status.Error(codes.Internal, "unable to verify request")
	}

	if !ok {
		return nil, status.Error(codes.Unauthenticated, "you cannot perform this action on this workspace")
	}

	logger.Debugw("granted permission", "workspace_id", res.Workspace.Id, "account_sub", acc.Sub)
	return ctx, nil
}

func ws(uc usecase.Controlplane, ctx context.Context, acc *authenticator.Account) (*usecase.WorkspaceGetByAccountRes, error) {
	wsId := metadata.Metadata(ctx).Get("x-kanthor-workspace")[0]
	if wsId == "" {
		return nil, errors.New("request without selected workspace is denied")
	}

	req := &usecase.WorkspaceGetByAccountReq{WorkspaceId: wsId, AccountSub: acc.Sub}
	res, err := uc.Workspace().GetByAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func action(method string) (string, string, error) {
	segments := strings.Split(method, "/")
	if len(segments) != 3 {
		return "", "", errors.New("malformed method")
	}

	return segments[1], segments[2], nil
}
