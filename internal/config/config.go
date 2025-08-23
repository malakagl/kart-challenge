package config

import (
	"fmt"
	"os"
	"time"

	validate "github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Logging    LoggingConfig    `yaml:"logging"`
	CouponCode CouponCodeConfig `yaml:"couponCode"`
}

type ServerConfig struct {
	Host                   string        `yaml:"host"`
	Port                   int           `yaml:"port" validate:"required"`
	MaxCouponCodeCacheSize int           `yaml:"maxCouponCodeCacheSize"`
	ReqLimitPerIPPerSec    int           `yaml:"reqLimitPerIPPerSec" validate:"min=1"`
	ReqBurstPerIPPerSec    int           `yaml:"reqBurstPerIPPerSec" validate:"min=1"`
	GracefulTimeout        time.Duration `yaml:"gracefulTimeout" validate:"required"`
}

type DatabaseConfig struct {
	Host                 string `yaml:"host" validate:"required"`
	Port                 int    `yaml:"port" validate:"required"`
	Name                 string `yaml:"name" validate:"required"`
	User                 string `yaml:"user" validate:"required"`
	Password             string `yaml:"password" validate:"required"`
	MigrationsFolderPath string `yaml:"migrationsFolderPath" validate:"required"`
	SSLMode              string `yaml:"sslMode"`
	Debug                bool   `yaml:"debug"`
	Type                 string `yaml:"type" validate:"required"` // e.g., "postgres", "sqlite"
}

type CouponCodeConfig struct {
	Unzipped  bool     `yaml:"unzipped"`
	FilePaths []string `yaml:"filePaths"`
}

type LoggingConfig struct {
	Level      string `yaml:"level" validate:"required"`
	JsonFormat bool   `yaml:"jsonFormat"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	if err := validate.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %v", err)
	}

	// add defaults
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}

	return cfg, nil
}
