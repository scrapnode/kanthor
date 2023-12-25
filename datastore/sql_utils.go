package datastore

import (
	"errors"
	"fmt"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/scrapnode/kanthor/pkg/timer"
	"gorm.io/gorm"
)

func SqlApplyScanQuery(tx *gorm.DB, timer timer.Timer, ns, prop string, query *entities.ScanningQuery) *gorm.DB {
	from := timer.UnixMilli(query.Start)
	to := timer.UnixMilli(query.End)
	low := suid.Id(ns, suid.BeforeTime(from))
	high := suid.Id(ns, suid.AfterTime(to))

	return tx.
		Where(fmt.Sprintf(`%s > ?`, prop), low).
		Where(fmt.Sprintf(`%s < ?`, prop), high).
		Limit(query.Limit)

}

func SqlError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}
	return err
}
