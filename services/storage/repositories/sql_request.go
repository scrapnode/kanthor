package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"gorm.io/gorm"
)

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) Create(ctx context.Context, docs []*entities.Request) ([]string, error) {

	if len(docs) == 0 {
		return []string{}, nil
	}

	datac := make(chan []string, 1)
	defer close(datac)

	errc := make(chan error, 1)
	defer close(errc)

	go func() {
		returning := []string{}

		names := []string{}
		values := map[string]interface{}{}
		for i := 0; i < len(docs); i++ {
			doc := docs[i]
			returning = append(returning, doc.Id)

			keys := []string{}
			for _, col := range entities.RequestProps {
				key := fmt.Sprintf("%s_%d", col, i)
				keys = append(keys, "@"+key)
				values[key] = entities.RequestMappers[col](doc)
			}

			names = append(names, fmt.Sprintf("(%s)", strings.Join(keys, ",")))
		}

		tableName := fmt.Sprintf(`"%s"`, entities.TableReq)
		columns := fmt.Sprintf(`"%s"`, strings.Join(entities.RequestProps, `","`))
		statement := fmt.Sprintf(
			"INSERT INTO %s(%s) VALUES %s ON CONFLICT(app_id, msg_id, id) DO NOTHING;",
			tableName,
			columns,
			strings.Join(names, ","),
		)

		if tx := sql.client.Exec(statement, values); tx.Error != nil {
			errc <- tx.Error
			return
		}

		datac <- returning
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case data := <-datac:
		return data, nil
	case err := <-errc:
		return nil, err
	}
}
