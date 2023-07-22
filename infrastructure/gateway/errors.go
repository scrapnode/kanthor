package gateway

import (
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Err400(msg string) error {
	return status.Error(codes.InvalidArgument, utils.Stringify(structure.M{"error": msg}))
}

func Err401(msg string) error {
	return status.Error(codes.Unauthenticated, utils.Stringify(structure.M{"error": msg}))
}

func Err403(msg string) error {
	return status.Error(codes.PermissionDenied, utils.Stringify(structure.M{"error": msg}))
}

func Err404(msg string) error {
	return status.Error(codes.NotFound, utils.Stringify(structure.M{"error": msg}))
}

func Err500(msg string) error {
	return status.Error(codes.Internal, utils.Stringify(structure.M{"error": msg}))
}
