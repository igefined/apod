package posgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const (
	dataSourceUrlFormat = "%s:%s@tcp(%s:%s)/%s?parseTime=true"
	dsnFormat           = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s"
)

func NewPostgresClient(cfg *postgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		dsnFormat,
		cfg.host,
		cfg.username,
		cfg.password,
		cfg.database,
		cfg.port,
		time.UTC,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres db: %v", err)
	}

	return db, err
}
