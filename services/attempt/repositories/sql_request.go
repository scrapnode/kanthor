package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"gorm.io/gorm"
)

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) Scan(ctx context.Context, appId string, msgIds []string, from, to time.Time) (map[string]Req, error) {
	if len(msgIds) == 0 {
		return map[string]Req{}, nil
	}

	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := entities.Id(entities.IdNsReq, suid.BeforeTime(from))
	high := entities.Id(entities.IdNsReq, suid.AfterTime(to))

	// @TODO: use chunk to fetch
	selects := []string{"app_id", "msg_id", "ep_id", "id", "tier"}
	var records []Req
	tx := sql.client.
		Table(entities.TableReq).
		Where("app_id = ?", appId).
		Where("msg_id IN ?", msgIds).
		Where("id > ?", low).
		Where("id < ?", high).
		Order("app_id DESC, msg_id DESC, id DESC").
		Select(selects).
		Find(&records)

	returning := map[string]Req{}
	for _, record := range records {
		returning[record.Id] = record
	}
	return returning, tx.Error
}

func (sql *SqlRequest) ListByIds(ctx context.Context, ids []string) ([]entities.Request, error) {
	rows, err := sql.client.
		Table(entities.TableMsg).
		Where("id IN ?", ids).
		Select([]string{"id", "timestamp", "msg_id", "ep_id", "tier", "app_id", "type", "metadata", "headers", "body", "uri", "method"}).
		Rows()

	if err != nil {
		return []entities.Request{}, err
	}

	var records []entities.Request
	defer rows.Close()
	for rows.Next() {
		record := entities.Request{}
		var metadata string
		var headers string
		var body string

		err := rows.Scan(
			&record.Id,
			&record.Timestamp,
			&record.MsgId,
			&record.EpId,
			&record.Tier,
			&record.AppId,
			&record.Type,
			&metadata,
			&headers,
			&body,
			&record.Uri,
			&record.Method,
		)

		// we don't accept partial success, if we got any error
		// return the error immediately
		if err != nil {
			return []entities.Request{}, err
		}
		if err := json.Unmarshal([]byte(metadata), &record.Metadata); err != nil {
			return []entities.Request{}, err
		}
		if err := json.Unmarshal([]byte(headers), &record.Headers); err != nil {
			return []entities.Request{}, err
		}
		record.Body = body
	}

	return records, nil
}
