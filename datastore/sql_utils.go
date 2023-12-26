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
