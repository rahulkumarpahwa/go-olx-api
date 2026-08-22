package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
	"github.com/rahulkumarpahwa/go-olx-api/internal/db"
	"github.com/rahulkumarpahwa/go-olx-api/internal/handlers"
	"github.com/rahulkumarpahwa/go-olx-api/internal/middleware"
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

	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("database connected...")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", handlers.Listings(DB))

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
	log.Printf("server is listening at http://localhost:%v\n", cfg.PORT)

	// graceful shutdown

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("%v\n", err.Error())
			log.Fatalf("Server Failed: %v\n", err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Println("server shutdown happens in...")
	for i := 5; i >= 0; i-- {
		log.Println(i)
		time.Sleep(time.Second * 1)
	}

	err = server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		log.Fatalf("Server Shutdown Failed: %v", err.Error())
	}
	log.Println("server shutdown gracefully.")
}
