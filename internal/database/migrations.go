package database

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	logging "github.com/malakagl/kart-challenge/pkg/logger"
)

func RunMigrations() error {
	dbURL := "postgres://user:password@localhost:5432/test?sslmode=disable"
	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		logging.Logger.Error().Msgf("Failed to create migration instance: %v", err)
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logging.Logger.Error().Msgf("Migration failed: %v", err)
		return err
	}

	return nil
}
