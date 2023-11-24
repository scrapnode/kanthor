package ds

import (
	"context"
	"encoding/json"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/project"
	"gorm.io/gorm"
)

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) Scan(ctx context.Context, appId string, msgIds []string, limit int) (map[string]Req, error) {
	if len(msgIds) == 0 {
		return map[string]Req{}, nil
	}

	selects := []string{"app_id", "msg_id", "ep_id", "id", "tier"}
	var requests []Req
	tx := sql.client.
		Table(entities.TableReq).
		Where("app_id = ?", appId).
		Where("msg_id IN ?", msgIds).
		Order("app_id ASC, msg_id ASC, id ASC").
		Limit(limit).
		Select(selects)

	if tx = tx.Find(&requests); tx.Error != nil {
		return nil, tx.Error
	}

	returning := map[string]Req{}
	// collect scanned records
	for _, s := range requests {
		returning[s.Id] = s
	}

	return returning, nil
}

func (sql *SqlRequest) ListByIds(ctx context.Context, ids []string) ([]entities.Request, error) {
	var returning []entities.Request
	for i := 0; i < len(ids); i += project.ScanBatchSize {
		j := utils.ChunkNext(i, len(ids), project.ScanBatchSize)

		requests, err := sql.list(ctx, ids[i:j])
		// we don't accept partial success, if we got any error
		// return the error immediately
		if err != nil {
			return nil, err
		}
		returning = append(returning, requests...)
	}

	return returning, nil
}

func (sql *SqlRequest) list(ctx context.Context, ids []string) ([]entities.Request, error) {
	rows, err := sql.client.
		Table(entities.TableReq).
		Where("id IN ?", ids).
		Select([]string{"id", "timestamp", "msg_id", "ep_id", "tier", "app_id", "type", "metadata", "headers", "body", "uri", "method"}).
		Rows()

	if err != nil {
		return []entities.Request{}, err
	}

	var records []entities.Request
	defer func() {
		if err := rows.Close(); err != nil {
			sql.client.Logger.Error(ctx, err.Error())
		}
	}()
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

		records = append(records, record)
	}

	return records, nil
}
