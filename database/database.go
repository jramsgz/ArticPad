package database

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Driver   string
	Host     string
	Username string
	Password string
	Port     int
	Database string
}

type Database struct {
	*gorm.DB
}

func New(config *DatabaseConfig) (*Database, error) {
	var db *gorm.DB
	var err error
	switch strings.ToLower(config.Driver) {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Database), &gorm.Config{})
		break
	case "postgresql", "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.Username, config.Password, config.Database)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		break
	}
	return &Database{db}, err
}
