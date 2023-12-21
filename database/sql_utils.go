package database

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func ApplyListQuery(tx *gorm.DB, q string, qcols []string, limit, page int) *gorm.DB {
	if len(q) >= 3 && len(qcols) > 0 {
		for _, qcol := range qcols {
			// because dataset volume of database is often small, so we can use scanning here
			tx = tx.Where(fmt.Sprintf(`"%s" LIKE ?`, qcol), "%"+q+"%")
		}
	}

	return tx.Limit(limit).Offset(MaxInt((page-1)*limit, 0))
}

func ApplyCountQuery(tx *gorm.DB, q string, qcols []string) *gorm.DB {
	if len(q) >= 3 && len(qcols) > 0 {
		for _, qcol := range qcols {
			// because dataset volume of database is often small, so we can use scanning here
			tx = tx.Where(fmt.Sprintf(`"%s" LIKE ?`, qcol), "%"+q+"%")
		}
	}

	return tx
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
