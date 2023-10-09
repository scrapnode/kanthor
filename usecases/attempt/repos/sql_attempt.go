package repos

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
)

type SqlAttempt struct {
	client *gorm.DB
}

var AttemptMapping = map[string]func(doc entities.Attempt) any{
	"req_id":           func(doc entities.Attempt) any { return doc.ReqId },
	"tier":             func(doc entities.Attempt) any { return doc.Tier },
	"status":           func(doc entities.Attempt) any { return doc.Status },
	"res_id":           func(doc entities.Attempt) any { return doc.ResId },
	"schedule_counter": func(doc entities.Attempt) any { return doc.ScheduleCounter },
	"schedule_next":    func(doc entities.Attempt) any { return doc.ScheduleNext },
	"scheduled_at":     func(doc entities.Attempt) any { return doc.ScheduledAt },
	"completed_at":     func(doc entities.Attempt) any { return doc.CompletedAt },
}
var AttemptMappingCols = lo.Keys(AttemptMapping)

func (sql *SqlAttempt) Create(ctx context.Context, docs []entities.Attempt) ([]string, error) {
	ids := []string{}

	if len(docs) == 0 {
		return ids, nil
	}

	names := []string{}
	values := map[string]interface{}{}
	for k, doc := range docs {
		ids = append(ids, doc.ReqId)

		keys := []string{}
		for _, col := range AttemptMappingCols {
			key := fmt.Sprintf("%s_%d", col, k)
			keys = append(keys, "@"+key)

			mapping := AttemptMapping[col]
			values[key] = mapping(doc)
		}
		names = append(names, fmt.Sprintf("(%s)", strings.Join(keys, ",")))
	}

	tableName := fmt.Sprintf(`"%s"`, (&entities.Attempt{}).TableName())
	columns := fmt.Sprintf(`"%s"`, strings.Join(AttemptMappingCols, `","`))
	statement := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES %s ON CONFLICT(id) DO NOTHING;",
		tableName,
		columns,
		strings.Join(names, ","),
	)

	if tx := sql.client.Exec(statement, values); tx.Error != nil {
		return nil, tx.Error
	}

	return ids, nil
}
