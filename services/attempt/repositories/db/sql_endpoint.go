package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlEndpoint struct {
	client *gorm.DB
}

func (sql *SqlEndpoint) Scan(ctx context.Context, query *entities.ScanningQuery) chan *entities.ScanningResult[[]entities.Endpoint] {
	ch := make(chan *entities.ScanningResult[[]entities.Endpoint], 1)
	go sql.scan(ctx, query, ch)
	return ch
}

func (sql *SqlEndpoint) scan(ctx context.Context, query *entities.ScanningQuery, ch chan *entities.ScanningResult[[]entities.Endpoint]) {
	defer close(ch)

	var cursor string
	for {
		if ctx.Err() != nil {
			return
		}

		tx := sql.client.
			Model(&entities.Endpoint{}).
			Order("id ASC").
			Limit(query.Size)
		if query.Search != "" {
			tx = tx.Where("id = ? ", query.Search)
		}
		if cursor != "" {
			tx = tx.Where("id < ?", cursor)
		}

		var data []entities.Endpoint
		if tx := tx.Find(&data); tx.Error != nil {
			ch <- &entities.ScanningResult[[]entities.Endpoint]{Error: tx.Error}
			return
		}

		ch <- &entities.ScanningResult[[]entities.Endpoint]{Data: data}

		if len(data) < query.Size {
			return
		}

		cursor = data[len(data)-1].Id
	}
}
