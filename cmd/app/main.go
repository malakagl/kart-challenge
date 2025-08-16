package main

import (
	"flag"
	"log"

	"github.com/malakagl/kart-challenge/internal/config"
	"github.com/malakagl/kart-challenge/internal/server"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg := config.LoadConfig(cfgPath)

	log.Println("Server setting up on port:", cfg.Server.Port)
	if err := server.Start(cfg); err != nil {
		log.Fatalf("server start up failed: %v", err)
	}
}
