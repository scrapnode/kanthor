package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
)

type SqlMessage struct {
	client *gorm.DB
}

func (sql *SqlMessage) Create(ctx context.Context, docs []entities.Message) ([]entities.TSEntity, error) {
	var records []entities.TSEntity
	if len(docs) == 0 {
		return records, nil
	}

	m := sql.mapper()
	if err := m.Parse(docs); err != nil {
		return nil, err
	}

	tableName := fmt.Sprintf(`"%s"`, entities.TableMsg)
	cols := fmt.Sprintf(`"%s"`, strings.Join(m.Names(), `","`))
	unnest := strings.Join(m.Casters(), ",")
	query := `INSERT INTO %s (%s) (SELECT * FROM UNNEST(%s)) ON CONFLICT (app_id, id) DO NOTHING;`
	statement := fmt.Sprintf(query, tableName, cols, unnest)

	if tx := sql.client.Exec(statement, m.Values()...); tx.Error != nil {
		return nil, tx.Error
	}

	for _, doc := range docs {
		record := entities.TSEntity{}
		record.Id = doc.Id
		record.Timestamp = doc.Timestamp
		records = append(records, record)
	}
	return records, nil
}

func (sql *SqlMessage) mapper() *datastore.Mapper[entities.Message] {
	return datastore.NewMapper[entities.Message](
		map[string]func(doc entities.Message) any{
			"id":        func(doc entities.Message) any { return doc.Id },
			"timestamp": func(doc entities.Message) any { return doc.Timestamp },
			"tier":      func(doc entities.Message) any { return doc.Tier },
			"app_id":    func(doc entities.Message) any { return doc.AppId },
			"type":      func(doc entities.Message) any { return doc.Type },
			"metadata":  func(doc entities.Message) any { return utils.Stringify(doc.Metadata) },
			"headers":   func(doc entities.Message) any { return utils.Stringify(doc.Headers) },
			"body":      func(doc entities.Message) any { return doc.Body },
		},
		map[string]string{
			// cast timestamp as int8[]
			"timestamp": "int8[]",
			// others will be varchar[] by default
		},
	)
}
