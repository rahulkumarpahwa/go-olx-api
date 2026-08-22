package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
	"github.com/rahulkumarpahwa/go-olx-api/internal/db"
	"github.com/rahulkumarpahwa/go-olx-api/internal/handlers"
	"github.com/rahulkumarpahwa/go-olx-api/internal/middleware"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("Config Loading Error: %v", err)
	}

	DB, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("Database Connection Error: %v", err)
	}

	defer DB.Close()

	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("database connected...")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)

	// logging every request
	loggedMux := middleware.LoggingMiddleware(mux)

	server := http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      loggedMux,
		ReadTimeout:  cfg.READTIMEOUT,
		WriteTimeout: cfg.WRITETIMEOUT,
		IdleTimeout:  cfg.IDLETIMEOUT,
	}

	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf("server is listening at http://localhost:%v", cfg.PORT)

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("%v", err.Error())
		log.Fatalf("Server Failed: %v", err.Error())
	}
}
