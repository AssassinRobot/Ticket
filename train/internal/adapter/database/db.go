package database

import (
	"fmt"
	"train/internal/application/core/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(url string) (*gorm.DB, error) {
	db, openErr := gorm.Open(postgres.Open((url)), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	if err := migrate(db, &domain.Train{}); err != nil {
		return nil, err
	}
	if err := migrate(db, &domain.Seat{}); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *gorm.DB, model any) error {
	if err := db.AutoMigrate(model); err != nil {
		return fmt.Errorf("db migration error: %v", err)
	}

	return nil
}
