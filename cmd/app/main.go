package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/malakagl/kart-challenge/internal/config"
	"github.com/malakagl/kart-challenge/internal/server"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg := config.LoadConfig(cfgPath)

	fmt.Println("Server starting on port:", cfg.Server.Port)
	if err := server.Start(cfg); err != nil {
		log.Fatalf("server start up failed: %v", err)
	}
}
