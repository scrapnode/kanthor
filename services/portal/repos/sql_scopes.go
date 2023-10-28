package repos

import (
	"fmt"

	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
)

func UseWsId(wsId, target string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		join := fmt.Sprintf(
			`JOIN "%s" ON "%s"."id" = "%s"."ws_id" AND "%s"."ws_id" = ?`,
			entities.TableWs,
			entities.TableWs,
			target,
			target,
		)
		return db.Joins(join, wsId)
	}
}

func UseAppId(appId, target string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		join := fmt.Sprintf(
			`JOIN "%s" ON "%s"."id" = "%s"."app_id" AND "%s"."app_id" = ?`,
			entities.TableApp,
			entities.TableApp,
			target,
			target,
		)
		return db.Joins(join, appId)
	}
}

func UseEpId(epId, target string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		join := fmt.Sprintf(
			`JOIN "%s" ON "%s"."id" = "%s"."ep_id" AND "%s"."ep_id" = ?`,
			entities.TableEp,
			entities.TableEp,
			target,
			target,
		)
		return db.Joins(join, epId)
	}
}
