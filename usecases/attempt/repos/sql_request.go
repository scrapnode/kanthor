package repos

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"gorm.io/gorm"
)

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) Scan(ctx context.Context, appId string, msgIds []string, from, to time.Time) (map[string]Req, error) {
	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := entities.Id(entities.IdNsReq, suid.BeforeTime(from))
	high := entities.Id(entities.IdNsReq, suid.AfterTime(to))

	selects := []string{"app_id", "msg_id", "ep_id", "id", "tier"}
	var records []Req
	tx := sql.client.
		Model(&entities.Request{}).
		Where("app_id = ?").
		Where("msg_id IN ?", msgIds).
		Where("id > ?", low).
		Where("id < ?", high).
		Order("app_id DESC, msg_id DESC, id DESC").
		Select(selects).
		Find(&records)

	returning := map[string]Req{}
	for _, record := range records {
		returning[record.Id] = record
	}
	return returning, tx.Error
}