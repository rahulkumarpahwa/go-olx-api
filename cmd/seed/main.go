package main

import (
	"log"

	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
	"github.com/rahulkumarpahwa/go-olx-api/internal/db"
	"github.com/rahulkumarpahwa/go-olx-api/seed"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("Config Loading Error: %v\n", err)
	}

	DB, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("Database Connection Error: %v\n", err)
	}
	defer DB.Close()

	if err := seed.SeedListings(DB); err != nil {
		log.Fatalf("Seeding Error: %v\n", err)
	}

	log.Println("database seeded with 10 dummy listings.")
}
