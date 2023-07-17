package database

import (
	"context"
	"github.com/scrapnode/kanthor/domain/structure"
	"gorm.io/gorm"
)

func SqlToListQuery(tx *gorm.DB, req structure.ListReq) *gorm.DB {
	if len(req.Ids) > 0 {
		tx = tx.Where("id IN ?", req.Ids)
	}

	tx = tx.Order("id DESC")

	if req.Limit > 0 {
		tx = tx.Limit(req.Limit)
	}

	if req.Cursor == "" {
		return tx
	}

	return tx.Where("id < ?", req.Cursor)
}

func SqlClientFromContext(ctx context.Context, client *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(CtxTransaction).(*gorm.DB); ok {
		return tx
	}
	return client
}
