package datastore

import (
	"errors"
	"fmt"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/identifier"
	"gorm.io/gorm"
)

type ScanningCondition struct {
	PrimaryKeyNs  string
	PrimaryKeyCol string
}

func SqlApplyScanQuery(tx *gorm.DB, query *entities.ScanningQuery, condition *ScanningCondition) *gorm.DB {
	low := identifier.Id(condition.PrimaryKeyNs, identifier.BeforeTime(query.From))
	high := identifier.Id(condition.PrimaryKeyNs, identifier.AfterTime(query.To))

	tx = tx.
		Where(fmt.Sprintf(`%s > ?`, condition.PrimaryKeyCol), low).
		Where(fmt.Sprintf(`%s < ?`, condition.PrimaryKeyCol), high).
		Limit(query.Limit)

	if query.Search != "" {
		// IMPORTANT: only support search by primary key
		// our primary key is often conbined from multiple columns
		// so you can search with the second column of the primary key
		// when and only when you added the first column to the where condition
		// for example, your primary key is message_pk(app_id, id)
		// you can only match the where condition for the ID column
		// when you add the where condition for the app_id column before
		// message: app_id = ? AND id = ?
		tx = tx.Where(fmt.Sprintf(`%s = ?`, condition.PrimaryKeyCol), query.Search)
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
