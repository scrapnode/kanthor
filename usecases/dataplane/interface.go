package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"net/http"
)

type Dataplane interface {
	patterns.Connectable
	Message() Message
}

type Message interface {
	Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error)
}

type MessagePutReq struct {
	AppId    string            `json:"app_id"`
	Type     string            `json:"type"`
	Headers  http.Header       `json:"headers"`
	Body     string            `json:"body"`
	Metadata map[string]string `json:"metadata"`
}

type MessagePutRes struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Bucket    string `json:"bucket"`
}
