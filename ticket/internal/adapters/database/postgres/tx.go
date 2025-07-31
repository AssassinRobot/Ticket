package postgres

import (
	"gorm.io/gorm"
)

type gormTx struct {
	db *gorm.DB
}

func begin(db *gorm.DB) (*gormTx, error) {
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &gormTx{db: tx}, nil
}

func (t *gormTx) Commit() error {
	return t.db.Commit().Error
}

func (t *gormTx) Rollback() error {
	return t.db.Rollback().Error
}
