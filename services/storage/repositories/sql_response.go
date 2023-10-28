package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
)

type SqlResponse struct {
	client *gorm.DB
}

var ResponseMapping = map[string]func(doc entities.Response) any{
	"id":        func(doc entities.Response) any { return doc.Id },
	"timestamp": func(doc entities.Response) any { return doc.Timestamp },
	"msg_id":    func(doc entities.Response) any { return doc.MsgId },
	"ep_id":     func(doc entities.Response) any { return doc.EpId },
	"req_id":    func(doc entities.Response) any { return doc.ReqId },
	"tier":      func(doc entities.Response) any { return doc.Tier },
	"app_id":    func(doc entities.Response) any { return doc.AppId },
	"type":      func(doc entities.Response) any { return doc.Type },
	"metadata":  func(doc entities.Response) any { return utils.Stringify(doc.Metadata) },
	"headers":   func(doc entities.Response) any { return utils.Stringify(doc.Headers) },
	"body":      func(doc entities.Response) any { return doc.Body },
	"uri":       func(doc entities.Response) any { return doc.Uri },
	"status":    func(doc entities.Response) any { return doc.Status },
	"error":     func(doc entities.Response) any { return doc.Error },
}
var ResponseMappingCols = lo.Keys(ResponseMapping)

func (sql *SqlResponse) Create(ctx context.Context, docs []entities.Response) ([]entities.TSEntity, error) {
	records := []entities.TSEntity{}

	if len(docs) == 0 {
		return records, nil
	}

	names := []string{}
	values := map[string]interface{}{}
	for k, doc := range docs {
		record := entities.TSEntity{}
		record.Id = doc.Id
		record.Timestamp = doc.Timestamp
		records = append(records, record)

		keys := []string{}
		for _, col := range ResponseMappingCols {
			key := fmt.Sprintf("%s_%d", col, k)
			keys = append(keys, "@"+key)

			mapping := ResponseMapping[col]
			values[key] = mapping(doc)
		}
		names = append(names, fmt.Sprintf("(%s)", strings.Join(keys, ",")))
	}

	tableName := fmt.Sprintf(`"%s"`, entities.TableRes)
	columns := fmt.Sprintf(`"%s"`, strings.Join(ResponseMappingCols, `","`))
	statement := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES %s ON CONFLICT(id) DO NOTHING;",
		tableName,
		columns,
		strings.Join(names, ","),
	)

	if tx := sql.client.Exec(statement, values); tx.Error != nil {
		return nil, tx.Error
	}

	return records, nil
}
