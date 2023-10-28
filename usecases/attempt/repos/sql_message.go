package repos

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

func (sql *SqlMessage) Scan(ctx context.Context, appId string, from, to time.Time, cursor string) ([]Msg, error) {
	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := entities.Id(entities.IdNsMsg, suid.BeforeTime(from))
	high := entities.Id(entities.IdNsMsg, suid.AfterTime(to))

	// @TODO: use chunk to fetch
	selects := []string{"app_id", "id", "tier", "timestamp"}
	var records []Msg
	tx := sql.client.
		Table(entities.TableMsg).
		Where("app_id = ?", appId).
		Where("id > ?", low).
		Where("id < ?", high).
		Order("app_id DESC, id DESC").
		Select(selects)

	if cursor != "" {
		tx = tx.Where("id", ">", cursor)
	}

	tx = tx.Find(&records)
	return records, tx.Error
}

func (sql *SqlMessage) ListByIds(ctx context.Context, ids []string) ([]entities.Message, error) {
	rows, err := sql.client.
		Table(entities.TableMsg).
		Where("id IN ?", ids).
		Select([]string{"id", "timestamp", "tier", "app_id", "type", "metadata", "headers", "body"}).
		Rows()

	if err != nil {
		return []entities.Message{}, err
	}

	var records []entities.Message
	defer rows.Close()
	for rows.Next() {
		record := entities.Message{}
		var metadata string
		var headers string
		var body string

		err := rows.Scan(
			&record.Id,
			&record.Timestamp,
			&record.Tier,
			&record.AppId,
			&record.Type,
			&metadata,
			&headers,
			&body,
		)

		// we don't accept partial success, if we got any error
		// return the error immediately
		if err != nil {
			return []entities.Message{}, err
		}
		if err := json.Unmarshal([]byte(metadata), &record.Metadata); err != nil {
			return []entities.Message{}, err
		}
		if err := json.Unmarshal([]byte(headers), &record.Headers); err != nil {
			return []entities.Message{}, err
		}
		record.Body = body
	}

	return records, nil
}
