package ds

import (
	"context"
)

type ResponseStatusRow struct {
	AppId  string `json:"app_id"`
	MsgId  string `json:"msg_id"`
	Id     string `json:"id"`
	EpId   string `json:"ep_id"`
	Status int    `json:"status"`
}

type Response interface {
	Scan(ctx context.Context, appId string, msgIds []string) (map[string]ResponseStatusRow, error)
}
