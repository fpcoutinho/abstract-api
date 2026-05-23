package main

import (
	"log"
	"os"

	"abstract-api/internal/server"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	engine := server.New()

	if err := engine.Run(addr); err != nil {
		log.Fatal(err)
	}
}