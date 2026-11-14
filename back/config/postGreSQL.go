package config

import (
	"net"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initPostgreSQL(cfg Config) (*gorm.DB, error) {
	connection := url.URL{
		Scheme: "postgres", User: url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host: net.JoinHostPort(cfg.DBHost, cfg.DBPort), Path: cfg.DBName,
	}
	params := url.Values{}
	if cfg.DBSSLMode != "" {
		params.Set("sslmode", cfg.DBSSLMode)
	}
	connection.RawQuery = params.Encode()

	db, err := gorm.Open(postgres.Open(connection.String()), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	return db, nil
}
