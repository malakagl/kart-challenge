package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/malakagl/kart-challenge/internal/config"
	"github.com/malakagl/kart-challenge/pkg/log"
)

func RunMigrations(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)
	log.Debug().Msgf("database migrations started: %v", dsn)
	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		log.Error().Msgf("Failed to create migration instance: %v", err)
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Error().Msgf("Migration failed: %v", err)
		return err
	}

	return nil
}
