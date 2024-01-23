package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/identifier"
	"gorm.io/gorm"
)

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) Scan(ctx context.Context, epId string, query *entities.ScanningQuery) chan *entities.ScanningResult[[]entities.Request] {
	ch := make(chan *entities.ScanningResult[[]entities.Request], 1)
	go sql.scan(ctx, epId, query, ch)
	return ch
}

func (sql *SqlRequest) scan(ctx context.Context, epId string, query *entities.ScanningQuery, ch chan *entities.ScanningResult[[]entities.Request]) {
	defer close(ch)

	low := identifier.Id(entities.IdNsReq, identifier.BeforeTime(query.From))
	high := identifier.Id(entities.IdNsReq, identifier.AfterTime(query.To))
	var cursor string
	for {
		if ctx.Err() != nil {
			return
		}

		tx := sql.client.
			Model(&entities.Request{}).
			Where("ep_id = ?", epId).
			Where("id > ?", low).
			Order("ep_id DESC, id DESC").
			Limit(query.Size)

		if query.Search != "" {
			tx = tx.Where("id = ?", query.Search)
		}

		if cursor == "" {
			tx = tx.Where("id < ?", high)
		} else {
			tx = tx.Where("id < ?", cursor)
		}

		var data []entities.Request
		if tx := tx.Find(&data); tx.Error != nil {
			ch <- &entities.ScanningResult[[]entities.Request]{Error: tx.Error}
			return
		}

		ch <- &entities.ScanningResult[[]entities.Request]{Data: data}

		if len(data) < query.Size {
			return
		}
	}
}
