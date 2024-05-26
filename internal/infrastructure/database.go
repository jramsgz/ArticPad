package infrastructure

import (
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type DatabaseConfig struct {
	Driver   string
	Host     string
	Username string
	Password string
	Port     int
	Database string
}

func connectToDB(config *DatabaseConfig) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	switch strings.ToLower(config.Driver) {
	case "sqlite":
		db, err = sqlx.Open("sqlite", config.Database)
		if err == nil {
			err = db.Ping()
		}
	case "postgresql", "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.Username, config.Password, config.Database)
		db, err = sqlx.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
		}
	default:
		return nil, fmt.Errorf("invalid database driver: %s", config.Driver)
	}
	return db, err
}
