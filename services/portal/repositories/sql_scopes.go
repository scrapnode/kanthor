package repositories

import (
	"fmt"

	"github.com/scrapnode/kanthor/internal/entities"
	"gorm.io/gorm"
)

func UseWsId(wsId, target string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		join := fmt.Sprintf(
			`JOIN "%s" ON "%s"."id" = "%s"."ws_id"`,
			entities.TableWs,
			entities.TableWs,
			target,
		)
		where := fmt.Sprintf(`"%s"."id" = ?`, entities.TableWs)

		return db.Joins(join).Where(where, wsId)
	}
}

func UseAppId(appId, target string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		join := fmt.Sprintf(
			`JOIN "%s" ON "%s"."id" = "%s"."app_id"`,
			entities.TableApp,
			entities.TableApp,
			target,
		)
		where := fmt.Sprintf(`"%s"."id" = ?`, entities.TableApp)

		return db.Joins(join).Where(where, appId)
	}
}

func UseEpId(epId, target string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		join := fmt.Sprintf(
			`JOIN "%s" ON "%s"."id" = "%s"."ep_id"`,
			entities.TableEp,
			entities.TableEp,
			target,
		)
		where := fmt.Sprintf(`"%s"."id" = ?`, entities.TableEp)

		return db.Joins(join).Where(where, epId)
	}
}
