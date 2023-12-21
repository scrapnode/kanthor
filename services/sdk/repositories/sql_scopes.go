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

func UseApp(target string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		join := fmt.Sprintf(
			`JOIN "%s" ON "%s"."id" = "%s"."app_id"`,
			entities.TableApp,
			entities.TableApp,
			target,
		)
		return db.Joins(join)
	}
}

func UseAppId(appId, target string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		where := fmt.Sprintf(`"%s"."id" = ?`, entities.IdNsApp)
		return UseApp(target)(db).Where(where, appId)
	}
}

func UseEp(target string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		join := fmt.Sprintf(
			`JOIN "%s" ON "%s"."id" = "%s"."ep_id"`,
			entities.TableEp,
			entities.TableEp,
			target,
		)
		return db.Joins(join)
	}
}

func UseEpId(epId, target string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		where := fmt.Sprintf(`"%s"."id" = ?`, entities.TableEp)
		return UseEp(target)(db).Where(where, epId)
	}
}
