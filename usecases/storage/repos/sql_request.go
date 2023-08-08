package repos

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
	"strings"
)

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) Create(ctx context.Context, docs []entities.Request) ([]entities.TSEntity, error) {
	records := []entities.TSEntity{}

	if len(docs) == 0 {
		return records, nil
	}

	cols := []string{"id", "timestamp", "bucket", "tier", "app_id", "type", "metadata", "headers", "body", "uri", "method"}
	names := []string{}
	values := map[string]interface{}{}
	for k, doc := range docs {
		record := entities.TSEntity{}
		record.Id = doc.Id
		record.Bucket = doc.Bucket
		record.Timestamp = doc.Timestamp
		records = append(records, record)

		keys := []string{}
		for _, col := range cols {
			key := fmt.Sprintf("%s_%d", col, k)
			keys = append(keys, "@"+key)
			switch col {
			case "id":
				values[key] = doc.Id
			case "timestamp":
				values[key] = doc.Timestamp
			case "bucket":
				values[key] = doc.Bucket
			case "tier":
				values[key] = doc.Tier
			case "app_id":
				values[key] = doc.AppId
			case "type":
				values[key] = doc.Type
			case "metadata":
				values[key] = utils.Stringify(doc.Metadata)
			case "headers":
				values[key] = utils.Stringify(doc.Headers)
			case "body":
				values[key] = string(doc.Body)
			case "uri":
				values[key] = doc.Uri
			case "method":
				values[key] = doc.Method
			default:
				return nil, fmt.Errorf("unknown column %s", col)
			}
		}
		names = append(names, fmt.Sprintf("(%s)", strings.Join(keys, ",")))
	}

	tableName := fmt.Sprintf(`"%s"`, (&entities.Request{}).TableName())
	columns := fmt.Sprintf(`"%s"`, strings.Join(cols, `","`))
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
