package infrastructure

import (
	"fmt"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Driver   string
	Host     string
	Username string
	Password string
	Port     int
	Database string
}

func connectToDB(config *DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	switch strings.ToLower(config.Driver) {
	case "sqlite":
		// TODO: Remove logger from production.
		db, err = gorm.Open(sqlite.Open(config.Database), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	case "postgresql", "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.Username, config.Password, config.Database)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
	default:
		return nil, fmt.Errorf("invalid database driver: %s", config.Driver)
	}
	return db, err
}
