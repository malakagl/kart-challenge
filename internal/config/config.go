package config

import (
	"log"
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
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port" validate:"required"`
	GracefulTimeout time.Duration `yaml:"gracefulTimeout" validate:"required"`
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
	Level      string `json:"level"`
	JsonFormat bool   `yaml:"jsonFormat"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("failed to read config file %s: %v\n", path, err)
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Printf("failed to parse config file: %v", err)
		return nil, err
	}

	if err := validate.New().Struct(cfg); err != nil {
		log.Printf("config validation failed: %v", err)
		return nil, err
	}

	return cfg, nil
}
