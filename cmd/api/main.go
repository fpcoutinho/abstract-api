package main

import (
	"log"
	"os"

	"abstract-api/internal/db"
	"abstract-api/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	// load .env for local development (no-op if file missing)
	_ = godotenv.Load()

	if err := db.InitFromEnv(); err != nil {
		log.Fatalf("db init: %v", err)
	}
	defer db.Close()

	engine := server.New()

	if err := engine.Run(addr); err != nil {
		log.Fatal(err)
	}
}
