package setup

import (
	"context"
	"errors"
	"fmt"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type demo struct {
	Workspace     entities.Workspace
	WorkspaceTier entities.WorkspaceTier
	Application   entities.Application
	Endpoint      entities.Endpoint
	EndpointRules []entities.EndpointRule
}

func Demo(conf *config.Config, logger logging.Logger, owner string, verbose bool) error {
	db := database.New(&conf.Database, logger)

	ctx := context.Background()
	if err := db.Connect(ctx); err != nil {
		return err
	}
	defer func() {
		if err := db.Disconnect(ctx); err != nil {
			logger.Error(err)
		}
	}()

	data := generate(owner)
	if verbose {
		showDemo(data)
	}

	// if we are using SQL database, it must be gorm
	client, ok := db.Client().(*gorm.DB)
	if ok {
		return client.Transaction(func(tx *gorm.DB) error {
			if wstx := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(data.Workspace); wstx.Error != nil {
				return wstx.Error
			}

			if wsttx := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(data.WorkspaceTier); wsttx.Error != nil {
				return wsttx.Error
			}

			if apptx := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(data.Application); apptx.Error != nil {
				return apptx.Error
			}

			if eptx := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(data.Endpoint); eptx.Error != nil {
				return eptx.Error
			}

			if eprtx := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(data.EndpointRules); eprtx.Error != nil {
				return eprtx.Error
			}

			return nil
		})
	}

	return errors.New("setup.demo: no demo data was written")
}

func generate(owner string) *demo {
	data := &demo{}

	data.Workspace = workspace(owner)
	data.WorkspaceTier = *data.Workspace.Tier
	// prevent gorm insert reference objects
	data.Workspace.Tier = nil
	data.Application = application(owner, data.Workspace)
	data.Endpoint = endpoint(owner, data.Application)
	data.EndpointRules = endpointRules(owner, data.Endpoint)

	return data
}

func workspace(owner string) entities.Workspace {
	ws := &entities.Workspace{
		OwnerId: owner,
		Name:    "demo",
	}
	ws.Id = id("ws", owner)
	ws.CreatedAt = time.Now().UTC().UnixMilli()
	ws.UpdatedAt = time.Now().UTC().UnixMilli()
	ws.Tier = &entities.WorkspaceTier{WorkspaceId: ws.Id, Name: "default"}
	return *ws
}

func application(owner string, ws entities.Workspace) entities.Application {
	app := &entities.Application{
		WorkspaceId: ws.Id,
		Name:        "demo",
	}
	app.Id = id("app", owner)
	app.CreatedAt = time.Now().UTC().UnixMilli()
	app.UpdatedAt = time.Now().UTC().UnixMilli()

	return *app
}

func endpoint(owner string, app entities.Application) entities.Endpoint {
	ep := &entities.Endpoint{
		AppId:  app.Id,
		Name:   "example.com",
		Method: "POST",
		Uri:    "https://example.com/",
	}
	ep.Id = id("ep", owner)
	ep.CreatedAt = time.Now().UTC().UnixMilli()
	ep.UpdatedAt = time.Now().UTC().UnixMilli()

	return *ep
}

func endpointRules(owner string, ep entities.Endpoint) []entities.EndpointRule {
	var items []entities.EndpointRule

	all := entities.EndpointRule{
		EndpointId:          ep.Id,
		Name:                "all",
		Priority:            0,
		Exclusionary:        false,
		ConditionSource:     "app_id",
		ConditionExpression: fmt.Sprintf("equal::%s", ep.AppId),
	}

	all.Id = id("epr", owner[0:utils.MinInt(8, len(owner))]+all.Name[0:utils.MinInt(8, len(all.Name))])
	all.CreatedAt = time.Now().UTC().UnixMilli()
	all.UpdatedAt = time.Now().UTC().UnixMilli()
	items = append(items, all)

	return items
}
