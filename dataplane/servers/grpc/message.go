package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/dataplane/servers/grpc/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MessageServer struct {
	protos.UnimplementedMessageServer
}

func (s *MessageServer) Put(context.Context, *protos.MessagePutReq) (*protos.MessagePutRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
