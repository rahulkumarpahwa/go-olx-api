package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rahulkumarpahwa/go-olx-api/internal/cfg"
)

func main() {
	cfg := cfg.Config()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello"))
	})

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(fmt.Sprintf(`{"message" : "Server is working fine", "status": "%v", "time" : "%v"}`, 200, time.DateTime)))
		if err != nil {

		}
	})

	server := http.Server{
		Addr:         cfg.PORT,
		Handler:      mux,
		ReadTimeout:  cfg.READTIMEOUT,
		WriteTimeout: cfg.WRITETIMEOUT,
		IdleTimeout:  cfg.IDLETIMEOUT,
	}

	fmt.Println("Welcome to the OLX API ")

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("%v", err.Error())
		log.Fatalf("Server Failed: %v", err.Error())
	}
}
