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

type SqlRequest struct {
	client *gorm.DB
}

var RequestMapping = map[string]func(doc entities.Request) any{
	"id":        func(doc entities.Request) any { return doc.Id },
	"timestamp": func(doc entities.Request) any { return doc.Timestamp },
	"msg_id":    func(doc entities.Request) any { return doc.MsgId },
	"ep_id":     func(doc entities.Request) any { return doc.EpId },
	"tier":      func(doc entities.Request) any { return doc.Tier },
	"app_id":    func(doc entities.Request) any { return doc.AppId },
	"type":      func(doc entities.Request) any { return doc.Type },
	"metadata":  func(doc entities.Request) any { return utils.Stringify(doc.Metadata) },
	"headers":   func(doc entities.Request) any { return utils.Stringify(doc.Headers) },
	"body":      func(doc entities.Request) any { return doc.Body },
	"uri":       func(doc entities.Request) any { return doc.Uri },
	"method":    func(doc entities.Request) any { return doc.Method },
}
var RequestMappingCols = lo.Keys(RequestMapping)

func (sql *SqlRequest) Create(ctx context.Context, docs []entities.Request) ([]entities.TSEntity, error) {
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
		for _, col := range RequestMappingCols {
			key := fmt.Sprintf("%s_%d", col, k)
			keys = append(keys, "@"+key)

			mapping := RequestMapping[col]
			values[key] = mapping(doc)
		}
		names = append(names, fmt.Sprintf("(%s)", strings.Join(keys, ",")))
	}

	tableName := fmt.Sprintf(`"%s"`, entities.TableReq)
	columns := fmt.Sprintf(`"%s"`, strings.Join(RequestMappingCols, `","`))
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
