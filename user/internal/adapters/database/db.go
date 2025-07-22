package database

import (
	"fmt"
	"user/internal/application/core/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(url string) (*gorm.DB, error) {
	db, openErr := gorm.Open(postgres.Open((url)), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	return db,nil
}
