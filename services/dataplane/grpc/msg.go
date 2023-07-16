package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/services/dataplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/dataplane"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type msg struct {
	protos.UnimplementedMsgServer
	service *dataplane
}

func (server *msg) Put(ctx context.Context, req *protos.MsgPutReq) (*protos.MsgPutRes, error) {
	request := &usecase.MessagePutReq{
		AppId:    req.AppId,
		Type:     req.Type,
		Headers:  http.Header{},
		Body:     req.Body,
		Metadata: map[string]string{},
	}
	for key, value := range req.Headers {
		request.Headers.Set(key, value)
	}
	for key, value := range req.Metadata {
		request.Metadata[key] = value
	}

	response, err := server.service.uc.Message().Put(ctx, request)
	if err != nil {
		server.service.logger.Error(err.Error(), "request", utils.Stringify(req))
		return nil, status.Error(codes.Internal, "oops, something went wrong")
	}

	res := &protos.MsgPutRes{Id: response.Id, Timestamp: response.Timestamp, Bucket: response.Bucket}
	return res, nil
}
