package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/services/dataplane/grpc/protos"
	usecase "github.com/scrapnode/kanthor/usecases/dataplane"
	"net/http"
)

type Message struct {
	protos.UnimplementedMessageServer
	service *dataplane
}

func (server *Message) Put(ctx context.Context, req *protos.PutReq) (*protos.PutRes, error) {
	server.service.meter.Counter("dataplane_message_put_total", 1)

	request := &usecase.PutMessageReq{
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

	response, err := server.service.uc.PutMessage(ctx, request)
	if err != nil {
		server.service.meter.Counter("dataplane_message_put_error", 1)
		server.service.logger.Error(err)
		return nil, err
	}

	res := &protos.PutRes{Id: response.Id, Timestamp: response.Timestamp, Bucket: response.Bucket}
	return res, nil
}
