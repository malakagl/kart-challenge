package config

import (
	"log"
	"os"

	validate "github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port" validate:"required"`
}

type DatabaseConfig struct {
	Host       string `json:"host" validate:"required"`
	Port       int    `json:"port" validate:"required"`
	Name       string `json:"name" validate:"required"`
	User       string `json:"user" validate:"required"`
	Password   string `json:"password" validate:"required"`
	SSLEnabled bool   `json:"sslEnabled"`
	Debug      bool   `json:"debug"`
	Type       string `json:"type" validate:"required"` // e.g., "postgres", "sqlite"
}

type LoggingConfig struct {
	Debug bool `json:"debug"`
}

func LoadConfig(path string) *Config {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	if err := validate.New().Struct(cfg); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	return cfg
}
