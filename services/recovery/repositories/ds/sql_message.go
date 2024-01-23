package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/identifier"
	"gorm.io/gorm"
)

type SqlMessage struct {
	client *gorm.DB
}

func (sql *SqlMessage) Scan(ctx context.Context, appId string, query *entities.ScanningQuery) chan *entities.ScanningResult[[]entities.Message] {
	ch := make(chan *entities.ScanningResult[[]entities.Message], 1)
	go sql.scan(ctx, appId, query, ch)
	return ch
}

func (sql *SqlMessage) scan(ctx context.Context, appId string, query *entities.ScanningQuery, ch chan *entities.ScanningResult[[]entities.Message]) {
	defer close(ch)

	low := identifier.Id(entities.IdNsMsg, identifier.BeforeTime(query.From))
	high := identifier.Id(entities.IdNsMsg, identifier.AfterTime(query.To))
	var cursor string
	for {
		if ctx.Err() != nil {
			return
		}

		tx := sql.client.
			Model(&entities.Message{}).
			Where("app_id = ?", appId).
			Where("id > ?", low).
			Order("app_id DESC, id DESC").
			Limit(query.Size)

		if query.Search != "" {
			tx = tx.Where("id = ?", query.Search)
		}

		if cursor == "" {
			tx = tx.Where("id < ?", high)
		} else {
			tx = tx.Where("id < ?", cursor)
		}

		var data []entities.Message
		if tx := tx.Find(&data); tx.Error != nil {
			ch <- &entities.ScanningResult[[]entities.Message]{Error: tx.Error}
			return
		}

		ch <- &entities.ScanningResult[[]entities.Message]{Data: data}

		if len(data) < query.Size {
			return
		}
	}
}
