package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/routing"
	"gorm.io/gorm"
)

type SqlApplication struct {
	client *gorm.DB
}

func (sql *SqlApplication) Get(ctx context.Context, id string) (*entities.Application, error) {
	doc := &entities.Application{}
	doc.Id = id
	if tx := sql.client.Model(doc).Where("id = ?", id).First(doc); tx.Error != nil {
		return nil, tx.Error
	}

	return doc, nil
}

func (sql *SqlApplication) GetRoutes(ctx context.Context, ids []string) (map[string][]routing.Route, error) {
	returning := make(map[string][]routing.Route)

	endpoints, err := sql.getRouteEndpoints(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(endpoints) == 0 {
		return returning, nil
	}

	rules, err := sql.getRouteRules(ctx, endpoints)
	if err != nil {
		return nil, err
	}
	if len(rules) == 0 {
		return returning, nil
	}

	for i := range endpoints {
		if _, has := returning[endpoints[i].AppId]; has {
			returning[endpoints[i].AppId] = append(returning[endpoints[i].AppId], routing.Route{
				Endpoint: &endpoints[i],
				Rules:    rules[endpoints[i].Id],
			})
			continue
		}

		returning[endpoints[i].AppId] = []routing.Route{{
			Endpoint: &endpoints[i],
			Rules:    rules[endpoints[i].Id],
		}}
	}

	return returning, nil
}

func (sql *SqlApplication) getRouteEndpoints(ctx context.Context, appIds []string) ([]entities.Endpoint, error) {
	var endpoints []entities.Endpoint
	tx := sql.client.Model(&entities.Endpoint{}).Where("app_id IN ?", appIds).Find(endpoints)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return endpoints, nil
}

func (sql *SqlApplication) getRouteRules(ctx context.Context, endpoints []entities.Endpoint) (map[string][]entities.EndpointRule, error) {
	returning := make(map[string][]entities.EndpointRule)
	var ids []string
	for i := range endpoints {
		returning[endpoints[i].Id] = make([]entities.EndpointRule, 0)
		ids = append(ids, endpoints[i].Id)
	}

	var rules []entities.EndpointRule
	tx := sql.client.
		Model(&entities.EndpointRule{}).
		Where("ep_id IN ?", ids).
		// IMPORTANT: we must get the exclusionary rule first to match it, then priority first
		Order("ep_id DESC, exclusionary DESC, priority DESC").
		Find(rules)
	if tx.Error != nil {
		return nil, tx.Error
	}

	for i := range rules {
		returning[rules[i].EpId] = append(returning[rules[i].EpId], rules[i])
	}
	return returning, nil
}
