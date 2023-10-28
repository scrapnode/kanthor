package repositories

import (
	"context"
	"time"
)

type Req struct {
	AppId string `json:"app_id"`
	MsgId string `json:"msg_id"`
	EpId  string `json:"ep_id"`
	Id    string `json:"id"`
	Tier  string `json:"tier"`
}

type Request interface {
	Scan(ctx context.Context, appId string, msgIds []string, from, to time.Time) (map[string]Req, error)
}
