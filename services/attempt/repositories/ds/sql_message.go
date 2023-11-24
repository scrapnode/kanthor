package ds

import (
	"context"
	"encoding/json"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"gorm.io/gorm"
)

type SqlMessage struct {
	client *gorm.DB
}

func (sql *SqlMessage) Count(ctx context.Context, appId string, from, to time.Time) (int64, error) {
	low := entities.Id(entities.IdNsMsg, suid.BeforeTime(from))
	high := entities.Id(entities.IdNsMsg, suid.AfterTime(to))

	var count int64
	tx := sql.client.
		Table(entities.TableMsg).
		Where("app_id = ?", appId).
		Where("id < ?", high).
		Where("id > ?", low).
		Count(&count)
	return count, tx.Error

}
func (sql *SqlMessage) Scan(ctx context.Context, appId string, from, to time.Time, limit int) chan *ScanResults[map[string]entities.Message] {
	ch := make(chan *ScanResults[map[string]entities.Message], 1)
	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := entities.Id(entities.IdNsMsg, suid.BeforeTime(from))
	high := entities.Id(entities.IdNsMsg, suid.AfterTime(to))

	go func() {
		var cursor string
		for {
			records := map[string]entities.Message{}

			tx := sql.client.
				Table(entities.TableMsg).
				Where("app_id = ?", appId).
				Where("id < ?", high).
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
				ch <- &ScanResults[map[string]entities.Message]{Error: err}

				close(ch)
				break
			}

			defer func() {
				if err := rows.Close(); err != nil {
					sql.client.Logger.Error(ctx, err.Error())
				}
			}()

			for rows.Next() {
				record := entities.Message{}
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
					ch <- &ScanResults[map[string]entities.Message]{Error: err}

					close(ch)
					return
				}

				if err := json.Unmarshal([]byte(metadata), &record.Metadata); err != nil {
					ch <- &ScanResults[map[string]entities.Message]{Error: err}

					close(ch)
					return
				}
				if err := json.Unmarshal([]byte(headers), &record.Headers); err != nil {
					ch <- &ScanResults[map[string]entities.Message]{Error: err}

					close(ch)
					return
				}

				records[record.Id] = record
				cursor = record.Id
			}

			ch <- &ScanResults[map[string]entities.Message]{Data: records}

			if len(records) < limit {
				close(ch)
				break
			}
		}
	}()

	return ch
}
