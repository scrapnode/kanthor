package repositories

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
}

func (sql *SqlApplication) Get(ctx context.Context, id string) (*entities.Application, error) {
	doc := &entities.Application{}
	doc.Id = id
	if tx := sql.client.Model(doc).Where("id = ?", id).First(doc); tx.Error != nil {
		return nil, fmt.Errorf("application.get: %w", tx.Error)
	}

	return doc, nil
}
