package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
	"github.com/rahulkumarpahwa/go-olx-api/internal/handlers"
	"github.com/rahulkumarpahwa/go-olx-api/internal/middleware"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("Config Loading Error: %v", err)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.Health)

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(fmt.Sprintf(`{"message" : "Server is working fine", "status": "%v", "time" : "%v"}`, 200, time.DateTime)))
		if err != nil {

		}
	})

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
	log.Printf("Server is Listening at http://localhost:%v", cfg.PORT)

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("%v", err.Error())
		log.Fatalf("Server Failed: %v", err.Error())
	}
}
