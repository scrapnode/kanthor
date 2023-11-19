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

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) Create(ctx context.Context, docs []entities.Request) ([]entities.TSEntity, error) {
	var records []entities.TSEntity
	if len(docs) == 0 {
		return records, nil
	}

	m := sql.mapper()
	if err := m.Parse(docs); err != nil {
		return nil, err
	}

	tableName := fmt.Sprintf(`"%s"`, entities.TableReq)
	cols := fmt.Sprintf(`"%s"`, strings.Join(m.Names(), `","`))
	unnest := strings.Join(m.Casters(), ",")
	query := `INSERT INTO %s (%s) (SELECT * FROM UNNEST(%s)) ON CONFLICT (app_id, msg_id, id) DO NOTHING;`
	statement := fmt.Sprintf(query, tableName, cols, unnest)

	if tx := sql.client.Exec(statement, m.Values()...); tx.Error != nil {
		return nil, tx.Error
	}

	return records, nil
}

func (sql *SqlRequest) mapper() *datastore.Mapper[entities.Request] {
	return datastore.NewMapper[entities.Request](
		map[string]func(doc entities.Request) any{
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
		},
		map[string]string{
			// cast timestamp as int8[]
			"timestamp": "int8[]",
			// others will be varchar[] by default
		},
	)
}
