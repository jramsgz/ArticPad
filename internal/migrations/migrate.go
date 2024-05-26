package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"github.com/rs/zerolog"
)

//go:embed sqlite/*.sql postgres/*.sql
var migrations embed.FS

// RunMigrations runs the migrations for the given database and dialect.
func RunMigrations(db *sql.DB, dialect string, logger zerolog.Logger) error {
	provider, err := newGooseProvider(db, dialect)
	if err != nil {
		return err
	}

	results, err := provider.Up(context.Background())
	if err != nil {
		return err
	}

	for _, result := range results {
		var prefix string
		if result.Error != nil {
			prefix = "ERR"
		} else if result.Empty {
			prefix = "-- "
		} else {
			prefix = "OK "
		}
		msg := fmt.Sprintf("%s %s (%s)", prefix, result.Source.Path, result.Duration)
		if result.Error != nil {
			logger.Error().Msgf("%s: %s", msg, result.Error)
		} else {
			logger.Info().Msg(msg)
		}
	}

	if len(results) == 0 {
		logger.Info().Msg("No migrations to run")
	} else {
		version, err := provider.GetDBVersion(context.Background())
		if err != nil {
			logger.Error().Msgf("Could not get DB version: %s", err)
		} else {
			logger.Info().Msgf("Migrated database to version: %d", version)
		}
	}

	return nil
}

// newGooseProvider creates a new goose provider for the given database and dialect.
// It makes use of the embedded migrations to find the migration files.
func newGooseProvider(db *sql.DB, dialect string) (*goose.Provider, error) {
	var dbDialect database.Dialect
	switch dialect {
	case "sqlite":
		dbDialect = database.DialectSQLite3
	case "postgres":
		dbDialect = database.DialectPostgres
	default:
		return nil, errors.New("dialect not supported")
	}

	dbMigrations, err := fs.Sub(migrations, dialect)
	if err != nil {
		return nil, err
	}

	return goose.NewProvider(
		dbDialect,
		db,
		dbMigrations,
	)
}
