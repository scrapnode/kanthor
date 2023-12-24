package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
)

func ApplyListQuery(tx *gorm.DB, props []string, search string, limit, page int) *gorm.DB {
	if len(search) >= 3 && len(props) > 0 {
		for _, qcol := range props {
			// because dataset volume of database is often small, so we can use scanning here
			tx = tx.Where(fmt.Sprintf(`"%s" LIKE ?`, qcol), "%"+search+"%")
		}
	}

	return tx.Limit(limit).Offset(utils.MaxInt((page-1)*limit, 0))
}

func ApplyCountQuery(tx *gorm.DB, props []string, search string) *gorm.DB {
	if len(search) >= 3 && len(props) > 0 {
		for _, qcol := range props {
			// because dataset volume of database is often small, so we can use scanning here
			tx = tx.Where(fmt.Sprintf(`"%s" LIKE ?`, qcol), "%"+search+"%")
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
