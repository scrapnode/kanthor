package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/dataplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/dataplane"
)

type Server struct {
	protos.UnimplementedDataplaneServer
	logger logging.Logger
	uc     usecase.Dataplane
}

func (server *Server) PutMessage(ctx context.Context, req *protos.PutMessageReq) (*protos.PutMessageRes, error) {
	request := &usecase.PutMessageReq{AppId: req.AppId, Type: req.Type, Body: req.Body}

	response, err := server.uc.PutMessage(ctx, request)
	if err != nil {
		server.logger.Error(err)
		return nil, err
	}

	res := &protos.PutMessageRes{Id: response.Id, Timestamp: response.Timestamp, Bucket: response.Bucket}
	return res, nil
}
