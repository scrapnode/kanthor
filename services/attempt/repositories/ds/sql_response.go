package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"gorm.io/gorm"
)

type SqlResponse struct {
	client *gorm.DB
}

func (sql *SqlResponse) Scan(ctx context.Context, appId string, msgIds []string) (map[string]ResponseStatusRow, error) {
	if len(msgIds) == 0 {
		return map[string]ResponseStatusRow{}, nil
	}

	records := map[string]ResponseStatusRow{}

	rows, err := sql.client.
		Table(entities.TableRes).
		Where("app_id = ? AND msg_id IN ?", appId, msgIds).
		// the order is important because it's not only sort as primary key order
		// but also use to only fetch the latest row of duplicated rows
		Order("app_id ASC, msg_id ASC, id ASC ").
		Select([]string{"app_id", "msg_id", "id", "ep_id", "status"}).
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
		record := ResponseStatusRow{}
		err := rows.Scan(
			&record.AppId,
			&record.MsgId,
			&record.Id,
			&record.EpId,
			&record.Status,
		)

		// we don't accept partial success, if we got any error
		// return the error immediately
		if err != nil {
			return nil, err
		}
		records[record.Id] = record
	}

	return records, nil
}
