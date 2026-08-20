package main

import (
	"errors"
	"flag"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
)

func main() {

	start := time.Now()

	flag.Parse()

	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("migration Config Env Error: %v", err)
	}

	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf("migrations started.")

	m, err := migrate.New("file://"+cfg.MIGRATIONFILESPATH, cfg.DATABASE_URL)

	if err != nil {
		log.Fatalf(`migration Error: %v`, err)
	}

	// checking the flag is up or down
	// up
	if len(flag.Args()) > 0 && flag.Args()[0] == "up" {
		if err := m.Up(); err != nil {

			// check if no migration is need then stop immediately
			if errors.Is(err, migrate.ErrNoChange) {
				log.Printf("no migrations to apply, took: %v seconds", time.Since(start).Seconds())
				return
			}

			log.Fatalf(`migration Up Error: %v`, err)
		}
	}

	// down
	if len(flag.Args()) > 0 && flag.Args()[0] == "down" {
		if err := m.Down(); err != nil {

			// check if no migration is need then stop immediately
			if errors.Is(err, migrate.ErrNoChange) {
				log.Printf("no migrations to apply, took: %v seconds", time.Since(start).Seconds())
				return
			}

			log.Fatalf(`migration Down Error: %v`, err)
		}
	}

	log.Printf("migrations completed, took: %v seconds", time.Since(start).Seconds())
}

// reference :  https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
