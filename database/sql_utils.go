package database

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
)

func SqlApplyListQuery(tx *gorm.DB, query *entities.PagingQuery, props []string) *gorm.DB {
	if len(query.Search) >= 3 && len(props) > 0 {
		for _, qcol := range props {
			// because dataset volume of database is often small, so we can use scanning here
			tx = tx.Where(fmt.Sprintf(`%s LIKE ?`, qcol), "%"+query.Search+"%")
		}
	}

	return tx.Limit(query.Limit).Offset(utils.Max((query.Page-1)*query.Limit, 0))
}

func SqlApplyCountQuery(tx *gorm.DB, query *entities.PagingQuery, props []string) *gorm.DB {
	if len(query.Search) >= 3 && len(props) > 0 {
		for _, qcol := range props {
			// because dataset volume of database is often small, so we can use scanning here
			tx = tx.Where(fmt.Sprintf(`%s LIKE ?`, qcol), "%"+query.Search+"%")
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
