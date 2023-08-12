package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
	"sort"
)

type SqlEndpointRule struct {
	client *gorm.DB
}

func (sql *SqlEndpointRule) List(ctx context.Context, epIds []string) ([]entities.EndpointRule, error) {
	docs := []entities.EndpointRule{}
	if len(epIds) == 0 {
		return docs, nil
	}

	tx := sql.client.WithContext(ctx).Where("endpoint_id in ?", epIds).Find(&docs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// descending order
	sort.SliceStable(docs, func(i, j int) bool {
		return docs[i].Priority > docs[j].Priority
	})

	return docs, nil
}
