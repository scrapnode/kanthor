package authz

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/domain/entities"
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

// DefaultProtected returns the list of method that is protected only by authentication not authorization
func DefaultProtected() map[string]bool {
	public := map[string]bool{
		"/kanthor.controlplane.v1.Account/ListWorkspaces": true,
		"/kanthor.controlplane.v1.Account/Get":            true,
	}

	return public
}
func UnaryServerInterceptor(
	logger logging.Logger,
	engine authorizator.Authorizator,
	protected map[string]bool,
) grpccore.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpccore.UnaryServerInfo,
		handler grpccore.UnaryHandler,
	) (resp interface{}, err error) {
		ctx, err = authorize(logger, engine, ctx, protected, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(
	logger logging.Logger,
	engine authorizator.Authorizator,
	protected map[string]bool,
) grpccore.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpccore.ServerStream,
		info *grpccore.StreamServerInfo,
		handler grpccore.StreamHandler,
	) error {
		ctx, err := authorize(logger, engine, ss.Context(), protected, info.FullMethod)
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
	protected map[string]bool,
	method string,
) (context.Context, error) {
	if public, ok := ctx.Value(gateway.AccessPublicable).(bool); ok && public {
		return ctx, nil
	}

	if protectable(method, protected) {
		return context.WithValue(ctx, gateway.AccessProtectable, true), nil
	}

	acc, ok := ctx.Value(authenticator.CtxAuthAccount).(*authenticator.Account)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "unknown account")
	}

	ws, ok := ctx.Value(authorizator.CtxAuthzWorkspace).(*entities.Workspace)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "unknown workspace")
	}

	// if authenticated account is owner of this workspace, should allow them does anything they want
	if acc.Sub == ws.OwnerId {
		logger.Debugw("owner permissions", "workspace_id", ws, "account_sub", acc.Sub)
		return ctx, nil
	}

	obj, act, err := action(method)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ok, err = engine.Enforce(acc.Sub, ws.Id, obj, act)
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

func protectable(value string, maps map[string]bool) bool {
	should, ok := maps[value]
	return ok && should
}
