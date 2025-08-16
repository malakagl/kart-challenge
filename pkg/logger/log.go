package logging

import (
	"os"
	"strings"
	"time"

	"github.com/malakagl/kart-challenge/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	Logger zerolog.Logger
)

func Init(serviceName string, cfg config.LoggingConfig) {
	zerolog.TimeFieldFormat = time.RFC3339

	levelStr := strings.ToLower(cfg.Level)
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}

	var output zerolog.Logger
	if !cfg.JsonFormat {
		output = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		output = zerolog.New(os.Stderr)
	}

	Logger = output.With().
		Timestamp().
		Str("service", serviceName).
		Logger().
		Level(level)

	log.Logger = Logger
}
