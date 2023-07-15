package authz

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/stream"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	grpccore "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func UnaryServerInterceptor(
	logger logging.Logger,
	engine authorizator.Authorizator,
) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		ctx, err = authorize(logger, engine, ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(
	logger logging.Logger,
	engine authorizator.Authorizator,
) grpccore.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpccore.ServerStream,
		info *grpccore.StreamServerInfo,
		handler grpccore.StreamHandler,
	) error {
		ctx, err := authorize(logger, engine, ss.Context(), info.FullMethod)
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
	engine authorizator.Authorizator,
	ctx context.Context,
	method string,
) (context.Context, error) {
	if public, ok := ctx.Value(gateway.AccessPublicable).(bool); ok && public {
		return ctx, nil
	}

	acc, ok := ctx.Value(authenticator.CtxAuthAccount).(*authenticator.Account)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "unknown account")
	}

	meta := gateway.ExtractIncoming(ctx)
	ws := meta.Get("x-kanthor-workspace")
	if ws == "" {
		return ctx, status.Error(codes.InvalidArgument, "unknown workspace")
	}

	obj, act, err := action(method)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ok, err = engine.Enforce(acc.Sub, ws, obj, act)
	if err != nil {
		logger.Errorw(err.Error(), "workspace_id", ws, "account_sub", acc.Sub)
		return nil, status.Error(codes.Internal, "unable to perform authorization")
	}

	if !ok {
		return nil, status.Error(codes.PermissionDenied, "you cannot perform the action on selected workspace")
	}

	logger.Debugw("granted permission", "workspace_id", ws, "account_sub", acc.Sub)
	return ctx, nil
}

func action(method string) (string, string, error) {
	segments := strings.Split(method, "/")
	if len(segments) != 3 {
		return "", "", errors.New("malformed method")
	}

	return segments[1], segments[2], nil
}
