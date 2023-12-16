package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/scrapnode/kanthor/internal/structure"
	"gorm.io/gorm"
)

func ApplyListQuery(tx *gorm.DB, limit, page int, q string, qcols []string) *gorm.DB {
	if len(q) >= 3 && len(qcols) > 0 {
		for _, qcol := range qcols {
			tx = tx.Where(fmt.Sprintf(`"%s" LIKE ?`, qcol), q)
		}
	}

	return tx.Limit(limit).Offset(MaxInt((page-1)*limit, 0))
}

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

func SqlTxnFromContext(ctx context.Context, client *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(CtxTransaction).(*gorm.DB); ok {
		return tx
	}
	return client
}

func SqlError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}
	return err
}

func MaxInt[T int | int32 | int64](x, y T) T {
	if x < y {
		return y
	}
	return x
}
