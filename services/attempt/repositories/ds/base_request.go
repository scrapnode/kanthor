package ds

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
)

type Req struct {
	AppId string `json:"app_id"`
	MsgId string `json:"msg_id"`
	EpId  string `json:"ep_id"`
	Id    string `json:"id"`
	Tier  string `json:"tier"`
}

type Request interface {
	Scan(ctx context.Context, appId string, msgIds []string, limit int) (map[string]Req, error)
	ListByIds(ctx context.Context, ids []string) ([]entities.Request, error)
}
