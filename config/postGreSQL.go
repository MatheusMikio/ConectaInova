package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initPostgreSQL(cfg Config) (*gorm.DB, error) {
	dbConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dbConnection), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
