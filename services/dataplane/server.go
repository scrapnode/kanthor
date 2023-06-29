package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/dataplane/grpc/protos"
	"github.com/scrapnode/kanthor/usecases"
)

type Server struct {
	protos.UnimplementedDataplaneServer
	logger  logging.Logger
	usecase usecases.Dataplane
}

func (server *Server) PutMessage(ctx context.Context, req *protos.PutMessageReq) (*protos.PutMessageRes, error) {
	request := &usecases.DataplanePutMessageReq{AppId: req.AppId, Type: req.Type, Body: req.Body}

	response, err := server.usecase.PutMessage(ctx, request)
	if err != nil {
		server.logger.Error(err)
		return nil, err
	}

	res := &protos.PutMessageRes{Id: response.Id, Timestamp: response.Timestamp, Bucket: response.Bucket}
	return res, nil
}
