package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Dataplane interface {
	patterns.Connectable
	PutMessage(ctx context.Context, req *PutMessageReq) (*PutMessageRes, error)
}

type PutMessageReq struct {
	AppId string `json:"app_id"`
	Type  string `json:"type"`
	Body  string `json:"body"`
}

type PutMessageRes struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Bucket    string `json:"bucket"`
}
