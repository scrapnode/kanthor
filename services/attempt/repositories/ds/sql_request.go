package ds

import (
	"context"
	"encoding/json"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/sourcegraph/conc/pool"
	"gorm.io/gorm"
)

type SqlRequest struct {
	client *gorm.DB
}

func (sql *SqlRequest) Scan(ctx context.Context, appId string, msgIds []string) (map[string]*entities.Request, error) {
	records := map[string]*entities.Request{}
	if len(msgIds) == 0 {
		return records, nil
	}

	seen := map[string]bool{}

	rows, err := sql.client.
		Table(entities.TableReq).
		Where("app_id = ? AND msg_id IN ?", appId, msgIds).
		// the order is important because it's not only sort as primary key order
		// but also use to only fetch the latest row of duplicated rows
		Order("app_id ASC, msg_id ASC, id ASC ").
		Select(entities.RequestProps).
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
		record := &entities.Request{}
		var metadata string
		var headers string

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
			&record.Body,
			&record.Uri,
			&record.Method,
		)
		// check duplicated requests
		key := ReqKey(record.MsgId, record.EpId)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = true

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

func (sql *SqlRequest) ListByIds(ctx context.Context, maps map[string]map[string][]string) (map[string]*entities.Request, error) {
	returning := safe.Map[*entities.Request]{}

	// IMPORTANT: we are not sure about how many goroutine can be spinned up here
	// so we have to use hard limit of the pool size here to prevent too many connection is created
	p := pool.New().WithMaxGoroutines(5)
	for key, childmaps := range maps {
		appId := key
		for childkey, values := range childmaps {
			msgId := childkey
			reqIds := values
			p.Go(func() {
				rows, err := sql.client.
					Table(entities.TableReq).
					Where("app_id = ? AND msg_id = ? AND id IN ?", appId, msgId, reqIds).
					Select([]string{"id", "timestamp", "msg_id", "ep_id", "tier", "app_id", "type", "metadata", "headers", "body", "uri", "method"}).
					Rows()

				if err != nil {
					sql.client.Logger.Error(ctx, err.Error())
					return
				}

				defer func() {
					if err := rows.Close(); err != nil {
						sql.client.Logger.Error(ctx, err.Error())
					}
				}()
				for rows.Next() {
					record := &entities.Request{}
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
						sql.client.Logger.Error(ctx, err.Error())
						return
					}
					if err := json.Unmarshal([]byte(metadata), &record.Metadata); err != nil {
						sql.client.Logger.Error(ctx, err.Error())
						return
					}
					if err := json.Unmarshal([]byte(headers), &record.Headers); err != nil {
						sql.client.Logger.Error(ctx, err.Error())
						return
					}
					record.Body = body

					returning.Set(record.Id, record)
				}
			})
		}
	}
	p.Wait()

	return returning.Data(), nil
}
