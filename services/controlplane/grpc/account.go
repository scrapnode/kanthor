package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/pipeline"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/transform"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

type account struct {
	protos.UnimplementedAccountServer
	service *controlplane
	pipe    pipeline.Middleware
}

func (server *account) Get(ctx context.Context, req *protos.AccountGetReq) (*protos.AccountEntity, error) {
	acc := ctx.Value(authenticator.CtxAuthAccount).(*authenticator.Account)
	return transform.Account(acc), nil
}

func (server *account) ListWorkspaces(ctx context.Context, req *protos.AccountListWorkspacesReq) (*protos.AccountListWorkspacesRes, error) {
	pipe := server.pipe(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, err = server.service.uc.Workspace().ListOfAccount(ctx, request.(*usecase.WorkspaceListOfAccountReq))
		return
	})

	acc := ctx.Value(authenticator.CtxAuthAccount).(*authenticator.Account)
	wsIds, err := server.service.authorizator.Tenants(acc.Sub)
	if err != nil {
		return nil, err
	}

	request := &usecase.WorkspaceListOfAccountReq{Account: acc, AssignedWorkspaceIds: wsIds}
	response, err := pipe(ctx, request)
	if err != nil {
		return nil, err
	}

	return transform.WorkspaceListOfAccountRes(ctx, response.(*usecase.WorkspaceListOfAccountRes)), nil
}
