package ds

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/status"
	"github.com/scrapnode/kanthor/pkg/suid"
	"gorm.io/gorm"
)

type SqlAttempt struct {
	client *gorm.DB
}

func (sql *SqlAttempt) Create(ctx context.Context, docs []*entities.Attempt) ([]string, error) {
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
			returning = append(returning, doc.ReqId)

			keys := []string{}
			for _, col := range entities.AttemptProps {
				key := fmt.Sprintf("%s_%d", col, i)
				keys = append(keys, "@"+key)
				values[key] = entities.AttemptMappers[col](doc)
			}

			names = append(names, fmt.Sprintf("(%s)", strings.Join(keys, ",")))
		}

		tableName := fmt.Sprintf(`"%s"`, entities.TableAtt)
		columns := fmt.Sprintf(`"%s"`, strings.Join(entities.AttemptProps, `","`))
		statement := fmt.Sprintf(
			"INSERT INTO %s(%s) VALUES %s ON CONFLICT(req_id) DO NOTHING;",
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

func (sql *SqlAttempt) Count(ctx context.Context, appId string, from, to time.Time, next int64) (int64, error) {
	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := suid.Id(entities.IdNsReq, suid.BeforeTime(from))
	high := suid.Id(entities.IdNsReq, suid.AfterTime(to))

	var count int64
	tx := sql.client.
		Table(entities.TableAtt).
		Where("req_id < ?", high).
		Where("req_id > ?", low).
		Where("schedule_next <= ?", next).
		Count(&count)
	return count, tx.Error
}

func (sql *SqlAttempt) Scan(ctx context.Context, from, to time.Time, next int64, limit int) chan *ScanResults[map[string]*entities.Attempt] {
	ch := make(chan *ScanResults[map[string]*entities.Attempt], 1)

	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := suid.Id(entities.IdNsReq, suid.BeforeTime(from))
	high := suid.Id(entities.IdNsReq, suid.AfterTime(to))

	go func() {
		defer close(ch)

		var cursor string
		for {
			records := map[string]*entities.Attempt{}

			tx := sql.client.
				Table(entities.TableAtt).
				Where("req_id < ?", high).
				Where("schedule_next <= ?", next).
				// the order is important because it's not only sort as primary key order
				// but also use to only fetch the latest row of duplicated rows
				Order("req_id ASC").
				Select(entities.AttemptProps).
				Limit(limit)

			if cursor == "" {
				tx = tx.Where("req_id > ?", low)
			} else {
				tx = tx.Where("req_id > ?", cursor)
			}

			rows, err := tx.Rows()
			if err != nil {
				ch <- &ScanResults[map[string]*entities.Attempt]{Error: err}
				return
			}

			defer func() {
				if err := rows.Close(); err != nil {
					sql.client.Logger.Error(ctx, err.Error())
				}
			}()

			for rows.Next() {
				record := &entities.Attempt{}

				err := rows.Scan(
					&record.ReqId,
					&record.MsgId,
					&record.AppId,
					&record.Tier,
					&record.Status,
					&record.ResId,
					&record.ScheduleCounter,
					&record.ScheduleNext,
					&record.ScheduledAt,
					&record.CompletedAt,
				)

				if err != nil {
					ch <- &ScanResults[map[string]*entities.Attempt]{Error: err}
					return
				}

				records[record.ReqId] = record
				// IMPORTANT: always update cursor
				cursor = record.ReqId
			}

			ch <- &ScanResults[map[string]*entities.Attempt]{Data: records}

			// if we found less than request size, that mean we were in last page
			if len(records) < limit {
				return
			}
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
