package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"gorm.io/gorm"
)

type SqlResponse struct {
	client *gorm.DB
}

func (sql *SqlResponse) Create(ctx context.Context, docs []entities.Response) ([]entities.TSEntity, error) {
	var records []entities.TSEntity
	if len(docs) == 0 {
		return records, nil
	}

	m := sql.mapper()
	if err := m.Parse(docs); err != nil {
		return nil, err
	}

	tableName := fmt.Sprintf(`"%s"`, entities.TableRes)
	cols := fmt.Sprintf(`"%s"`, strings.Join(m.Names(), `","`))
	unnest := strings.Join(m.Casters(), ",")
	query := `INSERT INTO %s (%s) (SELECT * FROM UNNEST(%s)) ON CONFLICT (app_id, msg_id, id) DO NOTHING;`
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

func (sql *SqlResponse) mapper() *datastore.Mapper[entities.Response] {
	return datastore.NewMapper[entities.Response](
		map[string]func(doc entities.Response) any{
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
		},
		map[string]string{
			// cast timestamp as int8[]
			"timestamp": "int8[]",
			// cast status as int2[]
			"status": "int2[]",
			// others will be varchar[] by default
		},
	)
}
