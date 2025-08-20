package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

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
		panic("failed to load config: " + err.Error())
	}

	logging.Init("kart-challenge", cfg.Logging)
	logging.Info().Msgf("Server setting up on port: %d", cfg.Server.Port)

	if err := server.Start(cfg); err != nil {
		logging.Logger.Fatal().Msgf("server failed to start: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigChan
	logging.Info().Msgf("Received signal: %s, shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		logging.Logger.Fatal().Msgf("Server shutdown failed: %v", err)
	}

	logging.Info().Msg("Server exited gracefully")
}
