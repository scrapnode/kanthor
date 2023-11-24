package ds

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/scrapnode/kanthor/project"
	"gorm.io/gorm"
)

type SqlResponse struct {
	client *gorm.DB
}

func (sql *SqlResponse) Scan(ctx context.Context, appId string, msgIds []string, from, to time.Time) (map[string]Res, error) {
	if len(msgIds) == 0 {
		return map[string]Res{}, nil
	}

	// convert timestamp to safe id, so we can the table efficiently with primary key
	low := entities.Id(entities.IdNsRes, suid.BeforeTime(from))
	high := entities.Id(entities.IdNsRes, suid.AfterTime(to))

	selects := []string{"app_id", "msg_id", "ep_id", "id", "tier", "req_id", "status"}

	var cursor string
	returning := map[string]Res{}

	for {
		var scanned []Res

		tx := sql.client.
			Table(entities.TableRes).
			Where("app_id = ?", appId).
			Where("msg_id IN ?", msgIds).
			Where("id < ?", high).
			Order("app_id ASC, msg_id ASC, id ASC").
			Limit(project.ScanBatchSize).
			Select(selects)

		if cursor == "" {
			tx = tx.Where("id > ?", low)
		} else {
			tx = tx.Where("id > ?", cursor)
		}

		if tx = tx.Find(&scanned); tx.Error != nil {
			return nil, tx.Error
		}

		// collect scanned records
		for _, s := range scanned {
			returning[s.Id] = s
		}

		// if we found less than request size, that mean we were in last page
		if len(scanned) < project.ScanBatchSize {
			break
		}

		cursor = scanned[len(scanned)-1].Id
	}

	if len(returning) == 0 {
		sql.client.Logger.Warn(ctx, "scanning return zero records", "from", low, "to", high)
	}

	return returning, nil
}
