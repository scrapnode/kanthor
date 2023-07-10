package metadata

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func Metadata(ctx context.Context) metadata.MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadata.Pairs()
	}
	return md
}
