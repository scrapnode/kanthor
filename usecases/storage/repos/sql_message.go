package repos

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
)

type SqlMessage struct {
	client *gorm.DB
}

var MessageMapping = map[string]func(doc entities.Message) any{
	"id":        func(doc entities.Message) any { return doc.Id },
	"timestamp": func(doc entities.Message) any { return doc.Timestamp },
	"tier":      func(doc entities.Message) any { return doc.Tier },
	"app_id":    func(doc entities.Message) any { return doc.AppId },
	"type":      func(doc entities.Message) any { return doc.Type },
	"metadata":  func(doc entities.Message) any { return utils.Stringify(doc.Metadata) },
	"headers":   func(doc entities.Message) any { return utils.Stringify(doc.Headers) },
	"body":      func(doc entities.Message) any { return doc.Body },
}
var MessageMappingCols = lo.Keys(MessageMapping)

func (sql *SqlMessage) Create(ctx context.Context, docs []entities.Message) ([]entities.TSEntity, error) {
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
		for _, col := range MessageMappingCols {
			key := fmt.Sprintf("%s_%d", col, k)
			keys = append(keys, "@"+key)

			mapping := MessageMapping[col]
			values[key] = mapping(doc)
		}
		names = append(names, fmt.Sprintf("(%s)", strings.Join(keys, ",")))
	}

	tableName := fmt.Sprintf(`"%s"`, entities.TableMsg)
	columns := fmt.Sprintf(`"%s"`, strings.Join(MessageMappingCols, `","`))
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
