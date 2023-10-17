package repos

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"gorm.io/gorm"
)

type SqlResponse struct {
	client *gorm.DB
}

func (sql *SqlResponse) Scan(ctx context.Context, appId string, msgIds []string, from, to time.Time) (map[string]Res, error) {
	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := entities.Id(entities.IdNsRes, suid.BeforeTime(from))
	high := entities.Id(entities.IdNsRes, suid.AfterTime(to))

	selects := []string{"app_id", "msg_id", "req_id", "id", "status"}
	var records []Res
	tx := sql.client.
		Model(&entities.Response{}).
		Where("app_id = ?").
		Where("msg_id IN ?", msgIds).
		Where("id > ?", low).
		Where("id < ?", high).
		Order("app_id DESC, msg_id DESC, id DESC").
		Select(selects).
		Find(&records)

	returning := map[string]Res{}
	for _, record := range records {
		returning[record.Id] = record
	}
	return returning, tx.Error
}