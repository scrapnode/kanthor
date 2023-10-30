package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
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

func (sql *SqlAttempt) BulkCreate(ctx context.Context, docs []entities.Attempt) ([]string, error) {
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

	tableName := fmt.Sprintf(`"%s"`, entities.TableAtt)
	columns := fmt.Sprintf(`"%s"`, strings.Join(AttemptMappingCols, `","`))
	statement := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES %s ON CONFLICT(req_id) DO NOTHING;",
		tableName,
		columns,
		strings.Join(names, ","),
	)

	if tx := sql.client.Exec(statement, values); tx.Error != nil {
		return nil, tx.Error
	}

	return ids, nil
}

func (sql *SqlAttempt) Scan(ctx context.Context, from, to time.Time, matching int64) ([]entities.Attempt, error) {
	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := entities.Id(entities.IdNsMsg, suid.BeforeTime(from))
	high := entities.Id(entities.IdNsMsg, suid.AfterTime(to))

	// @TODO: use chunk to fetch
	var records []entities.Attempt
	tx := sql.client.
		Table(entities.TableMsg).
		Where("id > ?", low).
		Where("id < ?", high).
		Where("schedule_next <= ?", matching).
		Order("req_id DESC").
		Find(&records)

	return records, tx.Error
}

func (sql *SqlAttempt) MarkComplete(ctx context.Context, reqId string, res *entities.Response) error {
	statement := fmt.Sprintf(
		"UPDATE %s SET completed_at = ?, status = ?, res_id = ? WHERE id = ?",
		entities.TableAtt,
	)

	if tx := sql.client.Exec(statement, res.Timestamp, res.Status, res.Id, reqId); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlAttempt) MarkReschedule(ctx context.Context, reqId string, ts int64) error {
	statement := fmt.Sprintf(
		"UPDATE %s SET schedule_counter = schedule_counter + 1, schedule_next = ?",
		entities.TableAtt,
	)

	if tx := sql.client.Exec(statement, ts, reqId); tx.Error != nil {
		return tx.Error
	}

	return nil
}
