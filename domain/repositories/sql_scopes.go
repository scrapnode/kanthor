package repositories

import (
	"fmt"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NotDeleted(timer timer.Timer, table schema.Tabler) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := fmt.Sprintf("%s.deleted_at == 0 OR %s.deleted_at > ?", table.TableName(), table.TableName())
		return db.Where(query, timer.Now().UnixMilli())
	}
}
