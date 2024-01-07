package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
}

func (sql *SqlApplication) Scan(ctx context.Context, query *entities.ScanningQuery) chan *entities.ScanningResult[[]entities.Application] {
	ch := make(chan *entities.ScanningResult[[]entities.Application], 1)
	go sql.scan(ctx, query, ch)
	return ch
}

func (sql *SqlApplication) scan(ctx context.Context, query *entities.ScanningQuery, ch chan *entities.ScanningResult[[]entities.Application]) {
	defer close(ch)

	var cursor string
	for {
		if ctx.Err() != nil {
			return
		}

		tx := sql.client.
			Table(entities.TableApp).
			Order("id ASC").
			Limit(query.Limit)
		if query.Search != "" {
			tx = tx.Where("id = ? ", query.Search)
		}
		if cursor != "" {
			tx = tx.Where("id < ?", cursor)
		}

		var data []entities.Application
		if tx := tx.Find(&data); tx.Error != nil {
			ch <- &entities.ScanningResult[[]entities.Application]{Error: tx.Error}
			return
		}

		ch <- &entities.ScanningResult[[]entities.Application]{Data: data}

		if len(data) < query.Limit {
			return
		}

		cursor = data[len(data)-1].Id
	}
}
