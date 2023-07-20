package grpc

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/gateway/grpc/protos"
)

type healthcheck struct {
	protos.UnimplementedHealthServer
}

func (server *healthcheck) Check(ctx context.Context, req *protos.HealthCheckRequest) (*protos.HealthCheckResponse, error) {
	res := &protos.HealthCheckResponse{
		Status: protos.HealthCheckResponse_SERVING,
	}
	return res, nil
}

func (server *healthcheck) Watch(req *protos.HealthCheckRequest, srv protos.Health_WatchServer) error {
	res := &protos.HealthCheckResponse{
		Status: protos.HealthCheckResponse_SERVING,
	}
	return srv.Send(res)
}
