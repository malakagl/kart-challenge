package main

import (
	"log"

	"github.com/MalakaGL/kart-challenge-lahiru/kart-challenge/internal/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatalf("server start up failed: %v", err)
	}
}
