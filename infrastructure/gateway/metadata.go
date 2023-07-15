package gateway

import (
	"context"
	"google.golang.org/grpc/metadata"
)

type MD metadata.MD

func ExtractIncoming(ctx context.Context) MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return MD(metadata.Pairs())
	}
	return MD(md)
}

func ExtractOutgoing(ctx context.Context) MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return MD(metadata.Pairs())
	}
	return MD(md)
}

func (m MD) Get(key string) string {
	value, ok := m[key]
	if !ok {
		return ""
	}
	return value[0]
}

func (m MD) Del(key string) MD {
	delete(m, key)
	return m
}

func (m MD) Set(key string, value string) MD {
	m[key] = []string{value}
	return m
}

func (m MD) Add(key string, value string) MD {
	m[key] = append(m[key], value)
	return m
}
