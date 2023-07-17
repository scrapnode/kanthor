package database

import (
	"errors"
	"gorm.io/gorm"
)

func ErrGet(tx *gorm.DB) error {
	if tx.Error == nil {
		return nil
	}

	// transform error
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	}

	return tx.Error
}
