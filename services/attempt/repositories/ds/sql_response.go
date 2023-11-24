package ds

import (
	"context"
	"encoding/json"

	"github.com/scrapnode/kanthor/domain/entities"
	"gorm.io/gorm"
)

type SqlResponse struct {
	client *gorm.DB
}

func (sql *SqlResponse) Scan(ctx context.Context, appId string, msgIds []string) (map[string]entities.Response, error) {
	if len(msgIds) == 0 {
		return map[string]entities.Response{}, nil
	}

	records := map[string]entities.Response{}

	rows, err := sql.client.
		Table(entities.TableRes).
		Where("app_id = ? AND msg_id IN ?", appId, msgIds).
		Select(entities.ResponseProps).
		Rows()

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			sql.client.Logger.Error(ctx, err.Error())
		}
	}()
	for rows.Next() {
		record := entities.Response{}
		var metadata string
		var headers string

		err := rows.Scan(
			&record.Id,
			&record.Timestamp,
			&record.MsgId,
			&record.EpId,
			&record.ReqId,
			&record.Tier,
			&record.AppId,
			&record.Type,
			&metadata,
			&headers,
			&record.Body,
			&record.Uri,
			&record.Status,
			&record.Error,
		)

		// we don't accept partial success, if we got any error
		// return the error immediately
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
