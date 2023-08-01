package repos

import (
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func UseWsId(target schema.Tabler, wsId string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		ws := &entities.Workspace{}
		join := fmt.Sprintf(
			`RIGHT JOIN "%s" ON "%s"."id" = "%s"."workspace_id"`,
			ws.TableName(),
			ws.TableName(),
			target.TableName(),
		)
		where := fmt.Sprintf(`"%s"."id" = ?`, ws.TableName())
		return db.Joins(join).Where(where, wsId)
	}
}

func JoinApp(target schema.Tabler) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		app := &entities.Application{}

		join := fmt.Sprintf(
			`RIGHT JOIN "%s" ON "%s"."id" = "%s"."app_id"`,
			app.TableName(),
			app.TableName(),
			target.TableName(),
		)
		return db.Joins(join)
	}
}

func JoinEp(target schema.Tabler) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		ep := &entities.Endpoint{}

		join := fmt.Sprintf(
			`RIGHT JOIN "%s" ON "%s"."id" = "%s"."endpoint_id"`,
			ep.TableName(),
			ep.TableName(),
			target.TableName(),
		)
		return db.Joins(join)
	}
}
