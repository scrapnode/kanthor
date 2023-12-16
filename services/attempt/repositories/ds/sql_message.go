package ds

import (
	"context"
	"encoding/json"
	"time"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"gorm.io/gorm"
)

type SqlMessage struct {
	client *gorm.DB
}

func (sql *SqlMessage) Count(ctx context.Context, appId string, from, to time.Time) (int64, error) {
	low := suid.Id(entities.IdNsMsg, suid.BeforeTime(from))
	high := suid.Id(entities.IdNsMsg, suid.AfterTime(to))

	var count int64
	tx := sql.client.
		Table(entities.TableMsg).
		Where("app_id = ?", appId).
		Where("id < ?", high).
		Where("id > ?", low).
		Count(&count)
	return count, tx.Error

}

func (sql *SqlMessage) Scan(ctx context.Context, appId string, from, to time.Time, limit int) chan *ScanResults[map[string]*entities.Message] {
	ch := make(chan *ScanResults[map[string]*entities.Message], 1)
	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := suid.Id(entities.IdNsMsg, suid.BeforeTime(from))
	high := suid.Id(entities.IdNsMsg, suid.AfterTime(to))

	go func() {
		defer close(ch)

		var cursor string
		for {
			records := map[string]*entities.Message{}

			tx := sql.client.
				Table(entities.TableMsg).
				Where("app_id = ?", appId).
				Where("id < ?", high).
				// the order is important because it's not only sort as primary key order
				// but also use to only fetch the latest row of duplicated rows
				Order("app_id ASC, id ASC").
				Select(entities.MessageProps).
				Limit(limit)

			if cursor == "" {
				tx = tx.Where("id > ?", low)
			} else {
				tx = tx.Where("id > ?", cursor)
			}

			rows, err := tx.Rows()
			if err != nil {
				ch <- &ScanResults[map[string]*entities.Message]{Error: err}
				return
			}

			defer func() {
				if err := rows.Close(); err != nil {
					sql.client.Logger.Error(ctx, err.Error())
				}
			}()

			for rows.Next() {
				record := &entities.Message{}
				var metadata string
				var headers string

				err := rows.Scan(
					&record.Id,
					&record.Timestamp,
					&record.Tier,
					&record.AppId,
					&record.Type,
					&metadata,
					&headers,
					&record.Body,
				)

				if err != nil {
					ch <- &ScanResults[map[string]*entities.Message]{Error: err}
					return
				}

				if err := json.Unmarshal([]byte(metadata), &record.Metadata); err != nil {
					ch <- &ScanResults[map[string]*entities.Message]{Error: err}
					return
				}
				if err := json.Unmarshal([]byte(headers), &record.Headers); err != nil {
					ch <- &ScanResults[map[string]*entities.Message]{Error: err}
					return
				}

				records[record.Id] = record
				// IMPORTANT: always update cursor
				cursor = record.Id
			}

			ch <- &ScanResults[map[string]*entities.Message]{Data: records}

			if len(records) < limit {
				return
			}
		}
	}()

	return ch
}

func (sql *SqlMessage) ListByIds(ctx context.Context, appId string, ids []string) (map[string]*entities.Message, error) {
	records := map[string]*entities.Message{}
	if len(ids) == 0 {
		return records, nil
	}

	tx := sql.client.
		Table(entities.TableMsg).
		Where("app_id = ?", appId).
		Where("id in ?", ids).
		// IMPORTANT: without this order, you may got weird behavior
		Order("app_id ASC, id ASC").
		Select(entities.MessageProps)

	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			sql.client.Logger.Error(ctx, err.Error())
		}
	}()

	for rows.Next() {
		record := &entities.Message{}
		var metadata string
		var headers string

		err := rows.Scan(
			&record.Id,
			&record.Timestamp,
			&record.Tier,
			&record.AppId,
			&record.Type,
			&metadata,
			&headers,
			&record.Body,
		)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(metadata), &record.Metadata); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(headers), &record.Headers); err != nil {
			return nil, err
		}

		records[record.Id] = record
	}

	return records, nil
}
