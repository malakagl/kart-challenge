package main

import (
	"log"

	"github.com/malakagl/kart-challenge/internal/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatalf("server start up failed: %v", err)
	}
}
