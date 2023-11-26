package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"gorm.io/gorm"
)

type SqlResponse struct {
	client *gorm.DB
}

func (sql *SqlResponse) Create(ctx context.Context, docs []entities.Response) ([]string, error) {
	returning := []string{}

	if len(docs) == 0 {
		return returning, nil
	}

	names := []string{}
	values := map[string]interface{}{}
	for i := 0; i < len(docs); i++ {
		doc := &docs[i]
		returning = append(returning, doc.Id)

		keys := []string{}
		for _, col := range entities.ResponseProps {
			key := fmt.Sprintf("%s_%d", col, i)
			keys = append(keys, "@"+key)
			values[key] = entities.ResponseMappers[col](doc)
		}

		names = append(names, fmt.Sprintf("(%s)", strings.Join(keys, ",")))
	}

	tableName := fmt.Sprintf(`"%s"`, entities.TableRes)
	columns := fmt.Sprintf(`"%s"`, strings.Join(entities.ResponseProps, `","`))
	statement := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES %s ON CONFLICT(app_id, msg_id, id) DO NOTHING;",
		tableName,
		columns,
		strings.Join(names, ","),
	)

	if tx := sql.client.Exec(statement, values); tx.Error != nil {
		return nil, tx.Error
	}

	return returning, nil
}
