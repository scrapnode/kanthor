package database

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/domain/structure"
	"gorm.io/gorm"
)

func SqlToListQuery(tx *gorm.DB, req *structure.ListReq, attr string) *gorm.DB {
	if len(req.Ids) > 0 {
		tx = tx.Where(fmt.Sprintf("%s IN ?", attr), req.Ids)
	}

	tx = tx.Order(fmt.Sprintf("%s DESC", attr))

	if req.Limit > 0 {
		tx = tx.Limit(req.Limit)
	}

	if req.Cursor == "" {
		return tx
	}

	return tx.Where(fmt.Sprintf("%s < ?", attr), req.Cursor)
}

func SqlClientFromContext(ctx context.Context, client *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(CtxTransaction).(*gorm.DB); ok {
		return tx
	}
	return client
}
