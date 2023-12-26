package datastore

import (
	"errors"
	"fmt"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/suid"
	"gorm.io/gorm"
)

func SqlApplyScanQuery(tx *gorm.DB, ns, prop string, query *entities.ScanningQuery) *gorm.DB {
	low := suid.Id(ns, suid.BeforeTime(query.From))
	high := suid.Id(ns, suid.AfterTime(query.To))

	tx = tx.
		Where(fmt.Sprintf(`%s > ?`, prop), low).
		Where(fmt.Sprintf(`%s < ?`, prop), high).
		Limit(query.Limit)

	if query.Search != "" {
		// IMPORTANT: only support search by primary key
		// our primary key is often conbined from multiple columns
		// so you need to add an condition that is matched with multiple columns
		// for example
		// message: app_id = ? AND id = ?
		tx = tx.Where(fmt.Sprintf(`%s = ?`, prop), query.Search)
	}

	return tx

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
