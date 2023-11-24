package ds

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
)

type SqlResponse struct {
	client *gorm.DB
}

func (sql *SqlResponse) Scan(ctx context.Context, appId string, msgIds []string, limit int) (map[string]Res, error) {
	if len(msgIds) == 0 {
		return map[string]Res{}, nil
	}

	selects := []string{"app_id", "msg_id", "ep_id", "id", "tier", "req_id", "status"}
	var responses []Res
	tx := sql.client.
		Table(entities.TableRes).
		Where("app_id = ?", appId).
		Where("msg_id IN ?", msgIds).
		Order("app_id ASC, msg_id ASC, id ASC").
		Limit(limit).
		Select(selects)

	if tx = tx.Find(&responses); tx.Error != nil {
		return nil, tx.Error
	}

	returning := map[string]Res{}
	// collect responses records
	for _, s := range responses {
		returning[s.Id] = s
	}

	return returning, nil
}
