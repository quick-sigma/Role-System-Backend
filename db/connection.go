package db

import (
	"fmt"

	"motor-de-rol/backend/domain"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Connect(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Character{},
		&domain.Race{},
		&domain.CharacterStat{},
	)
}
