package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/protos"
)

type ws struct {
	protos.UnimplementedWsServer
	service *controlplane
}

func (server *ws) List(ctx context.Context, req *protos.WsListReq) (*protos.WsListRes, error) {
	res := &protos.WsListRes{}
	return res, nil
}

func (server *ws) Get(ctx context.Context, req *protos.WsGetReq) (*protos.Workspace, error) {
	res := &protos.Workspace{}
	return res, nil
}
