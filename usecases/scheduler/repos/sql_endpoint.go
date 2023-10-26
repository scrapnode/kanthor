package repos

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
)

type SqlEndpoint struct {
	client *gorm.DB
}

func (sql *SqlEndpoint) List(ctx context.Context, appId string) ([]entities.Endpoint, error) {
	docs := []entities.Endpoint{}
	tx := sql.client.WithContext(ctx).Where("app_id = ?", appId).Find(&docs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return docs, nil
}

// Rules return an ordered rules slice, exclusionary first then large priority
func (sql *SqlEndpoint) Rules(ctx context.Context, appId string) ([]entities.EndpointRule, error) {
	rules := []entities.EndpointRule{}

	ep := &entities.Endpoint{}
	epr := &entities.EndpointRule{}
	join := fmt.Sprintf("JOIN %s ON %s.id = %s.ep_id AND %s.app_id = ?", ep.TableName(), ep.TableName(), epr.TableName(), ep.TableName())
	order := fmt.Sprintf("%s.exclusionary DESC, %s.priority", epr.TableName(), epr.TableName())
	selects := fmt.Sprintf("%s.*", epr.TableName())

	tx := sql.client.WithContext(ctx).
		Joins(join, appId).
		Order(order).
		Select(selects).
		Find(&rules)

	return rules, tx.Error
}
