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
	log.SetFlags(log.Ldate | log.Ltime)

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatalf("invalid args. \n usage: go run < migrate path > < up | down > ")
	}

	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("migration Config Env Error: %v", err)
	}

	log.Printf("migrations started.")

	m, err := migrate.New("file://"+cfg.MIGRATIONFILESPATH, cfg.DATABASE_URL)
	if err != nil {
		log.Fatalf(`migration Error: %v`, err)
	}

	// checking the flag is up or down

	switch flag.Args()[0] {
	case "up":
		if err := m.Up(); err != nil {

			// check if no migration is need then stop immediately
			if errors.Is(err, migrate.ErrNoChange) {
				log.Printf("no '%s' migrations to apply, took: %.2f seconds", flag.Args()[0], time.Since(start).Seconds())
				return
			}

			log.Fatalf(`migration Up Error: %v`, err)
		}

	case "down":
		if err := m.Down(); err != nil {

			// check if no migration is need then stop immediately
			if errors.Is(err, migrate.ErrNoChange) {
				log.Printf("no '%s' migrations to apply, took: %.2f seconds", flag.Args()[0], time.Since(start).Seconds())
				return
			}

			log.Fatalf(`migration Down Error: %v`, err)
		}
	default:
		log.Fatalf("unknown command: %s, took: %.2f seconds", flag.Args()[0], time.Since(start).Seconds())
	}

	log.Printf("migrations '%s' completed, took: %.2f seconds", flag.Args()[0], time.Since(start).Seconds())
}

// reference :  https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
