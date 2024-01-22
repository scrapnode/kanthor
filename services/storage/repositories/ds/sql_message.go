package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SqlMessage struct {
	client *gorm.DB
}

func (sql *SqlMessage) Create(ctx context.Context, docs []*entities.Message) ([]string, error) {
	if len(docs) == 0 {
		return []string{}, nil
	}

	datac := make(chan []string, 1)
	defer close(datac)

	errc := make(chan error, 1)
	defer close(errc)

	go func() {
		tx := sql.client.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&docs)
		if tx.Error != nil {
			errc <- tx.Error
			return
		}

		returning := []string{}
		for i := range docs {
			returning = append(returning, docs[i].Id)
		}

		datac <- returning
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case data := <-datac:
		return data, nil
	case err := <-errc:
		return nil, err
	}
}
