package db

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
