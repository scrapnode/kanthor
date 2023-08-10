package repos

import (
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func UseWsId(wsId string, target schema.Tabler) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		ws := &entities.Workspace{}
		join := fmt.Sprintf(
			`JOIN "%s" ON "%s"."id" = "%s"."workspace_id" AND "%s"."workspace_id" = ?`,
			ws.TableName(),
			ws.TableName(),
			target.TableName(),
			target.TableName(),
		)
		return db.Joins(join, wsId)
	}
}

func UseAppId(appId string, target schema.Tabler) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		app := &entities.Application{}

		join := fmt.Sprintf(
			`JOIN "%s" ON "%s"."id" = "%s"."app_id" AND "%s"."app_id" = ?`,
			app.TableName(),
			app.TableName(),
			target.TableName(),
			target.TableName(),
		)
		return db.Joins(join, appId)
	}
}

func UseEpId(epId string, target schema.Tabler) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		ep := &entities.Endpoint{}

		join := fmt.Sprintf(
			`JOIN "%s" ON "%s"."id" = "%s"."endpoint_id" AND "%s"."endpoint_id" = ?`,
			ep.TableName(),
			ep.TableName(),
			target.TableName(),
			target.TableName(),
		)
		return db.Joins(join, epId)
	}
}
