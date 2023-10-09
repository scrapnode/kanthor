package repos

import (
	"context"
	"time"
)

type Res struct {
	AppId  string `json:"app_id"`
	MsgId  string `json:"msg_id"`
	EpId   string `json:"ep_id"`
	Id     string `json:"id"`
	Status int    `json:"status"`
}

type Response interface {
	Scan(ctx context.Context, appId string, msgIds []string, from, to time.Time) (map[string]Res, error)
}
