package main

import (
	"flag"
	"log"

	"github.com/malakagl/kart-challenge/internal/config"
	"github.com/malakagl/kart-challenge/internal/server"
	logging "github.com/malakagl/kart-challenge/pkg/log"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logging.Init("kart-challenge", cfg.Logging)
	logging.Info().Msgf("Server setting up on port: %d", cfg.Server.Port)
	if err := server.Start(cfg); err != nil {
		logging.Logger.Fatal().Msgf("server start up failed: %v", err)
	}
}
