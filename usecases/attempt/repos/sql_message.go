package repos

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"gorm.io/gorm"
)

type SqlMessage struct {
	client *gorm.DB
}

func (sql *SqlMessage) Scan(ctx context.Context, appId string, from, to time.Time, cursor string) ([]Msg, error) {
	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := entities.Id(entities.IdNsMsg, suid.BeforeTime(from))
	high := entities.Id(entities.IdNsMsg, suid.AfterTime(to))

	// @TODO: use chunk to fetch
	selects := []string{"app_id", "id", "tier", "timestamp"}
	var records []Msg
	tx := sql.client.
		Table((&entities.Message{}).TableName()).
		Where("app_id = ?", appId).
		Where("id > ?", low).
		Where("id < ?", high).
		Order("app_id DESC, id DESC").
		Select(selects)

	if cursor != "" {
		tx = tx.Where("id", ">", cursor)
	}

	tx = tx.Find(&records)
	return records, tx.Error
}

func (sql *SqlMessage) ListByIds(ctx context.Context, ids []string) ([]entities.Message, error) {
	var records []entities.Message

	tx := sql.client.
		Model(&entities.Message{}).
		Where("id IN ?", ids).
		Find(&records)

	return records, tx.Error
}
