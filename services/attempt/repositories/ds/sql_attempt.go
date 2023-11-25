package ds

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/scrapnode/kanthor/datastore"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/status"
	"github.com/scrapnode/kanthor/pkg/suid"
	"gorm.io/gorm"
)

type SqlAttempt struct {
	client *gorm.DB
}

func (sql *SqlAttempt) Create(ctx context.Context, docs []entities.Attempt) ([]string, error) {
	var ids []string
	if len(docs) == 0 {
		return ids, nil
	}

	m := sql.mapper()
	if err := m.Parse(docs); err != nil {
		return nil, err
	}

	tableName := fmt.Sprintf(`"%s"`, entities.TableAtt)
	cols := fmt.Sprintf(`"%s"`, strings.Join(m.Names(), `","`))
	unnest := strings.Join(m.Casters(), ",")
	query := `INSERT INTO %s (%s) (SELECT * FROM UNNEST(%s)) ON CONFLICT (req_id) DO NOTHING;`
	statement := fmt.Sprintf(query, tableName, cols, unnest)

	if tx := sql.client.Exec(statement, m.Values()...); tx.Error != nil {
		return nil, tx.Error
	}

	for _, doc := range docs {
		ids = append(ids, doc.ReqId)
	}
	return ids, nil
}

func (sql *SqlAttempt) mapper() *datastore.Mapper[entities.Attempt] {
	return datastore.NewMapper[entities.Attempt](
		map[string]func(doc entities.Attempt) any{
			"req_id":           func(doc entities.Attempt) any { return doc.ReqId },
			"msg_id":           func(doc entities.Attempt) any { return doc.MsgId },
			"app_id":           func(doc entities.Attempt) any { return doc.AppId },
			"tier":             func(doc entities.Attempt) any { return doc.Tier },
			"status":           func(doc entities.Attempt) any { return doc.Status },
			"res_id":           func(doc entities.Attempt) any { return doc.ResId },
			"schedule_counter": func(doc entities.Attempt) any { return doc.ScheduleCounter },
			"schedule_next":    func(doc entities.Attempt) any { return doc.ScheduleNext },
			"scheduled_at":     func(doc entities.Attempt) any { return doc.ScheduledAt },
			"completed_at":     func(doc entities.Attempt) any { return doc.CompletedAt },
		},
		map[string]string{
			// cast timestamp as int8[]
			"timestamp": "int8[]",
			// cast status as int2[]
			"status": "int2[]",
			// cast schedule_counter as int2[]
			"schedule_counter": "int2[]",
			// cast schedule_next as int8[]
			"schedule_next": "int8[]",
			// cast scheduled_at as int8[]
			"scheduled_at": "int8[]",
			// cast completed_at as int8[]
			"completed_at": "int8[]",
			// others will be varchar[] by default
		},
	)
}

func (sql *SqlAttempt) Count(ctx context.Context, appId string, from, to time.Time, next int64) (int64, error) {
	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := entities.Id(entities.IdNsReq, suid.BeforeTime(from))
	high := entities.Id(entities.IdNsReq, suid.AfterTime(to))

	var count int64
	tx := sql.client.
		Table(entities.TableAtt).
		Where("req_id < ?", high).
		Where("req_id > ?", low).
		Where("schedule_next <= ?", next).
		Count(&count)
	return count, tx.Error
}

func (sql *SqlAttempt) Scan(ctx context.Context, from, to time.Time, next int64, limit int) chan *ScanResults[[]entities.Attempt] {
	ch := make(chan *ScanResults[[]entities.Attempt], 1)

	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := entities.Id(entities.IdNsReq, suid.BeforeTime(from))
	high := entities.Id(entities.IdNsReq, suid.AfterTime(to))

	go func() {
		defer close(ch)

		var cursor string
		for {
			var records []entities.Attempt

			tx := sql.client.
				Table(entities.TableAtt).
				Where("req_id < ?", high).
				Where("schedule_next <= ?", next).
				Order("req_id ASC").
				Limit(limit)

			if cursor == "" {
				tx = tx.Where("req_id > ?", low)
			} else {
				tx = tx.Where("req_id > ?", cursor)
			}

			if tx = tx.Find(&records); tx.Error != nil {

				ch <- &ScanResults[[]entities.Attempt]{Error: tx.Error}
				return
			}

			ch <- &ScanResults[[]entities.Attempt]{Data: records}

			// if we found less than request size, that mean we were in last page
			if len(records) < limit {
				return
			}

			// IMPORTANT: always update cursor
			cursor = records[len(records)-1].ReqId
		}
	}()

	return ch
}

func (sql *SqlAttempt) MarkComplete(ctx context.Context, reqId string, res *entities.Response) error {
	statement := fmt.Sprintf(
		"UPDATE %s SET completed_at = ?, status = ?, res_id = ? WHERE req_id = ?",
		entities.TableAtt,
	)

	if tx := sql.client.Exec(statement, res.Timestamp, res.Status, res.Id, reqId); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlAttempt) MarkReschedule(ctx context.Context, reqId string, ts int64) error {
	statement := fmt.Sprintf(
		"UPDATE %s SET schedule_counter = schedule_counter + 1, schedule_next = ? WHERE req_id = ?",
		entities.TableAtt,
	)

	if tx := sql.client.Exec(statement, ts, reqId); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (sql *SqlAttempt) MarkIgnore(ctx context.Context, reqIds []string) error {
	statement := fmt.Sprintf(
		"UPDATE %s SET status = ? WHERE req_id IN ?",
		entities.TableAtt,
	)

	if tx := sql.client.Exec(statement, status.ErrIgnore, reqIds); tx.Error != nil {
		return tx.Error
	}

	return nil
}
