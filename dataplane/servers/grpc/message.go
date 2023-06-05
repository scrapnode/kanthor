package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/dataplane/servers/grpc/protos"
	"github.com/scrapnode/kanthor/dataplane/services"
)

type MessageServer struct {
	protos.UnimplementedMessageServer
	service services.Message
}

func (s *MessageServer) Put(ctx context.Context, req *protos.MessageCreateReq) (*protos.MessageCreateRes, error) {
	request := &services.MessageCreateReq{AppId: req.AppId, Type: req.Type, Body: req.Body}

	response, err := s.service.Create(ctx, request)
	if err != nil {
		return nil, err
	}

	res := &protos.MessageCreateRes{Id: response.Id, Timestamp: response.Timestamp.Unix(), Bucket: response.Bucket}
	return res, nil
}
