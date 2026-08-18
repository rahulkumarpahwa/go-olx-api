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
	server := http.NewServeMux()

	server.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Hello"))
	})

	server.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(200)
		_, err := w.Write([]byte(fmt.Sprintf(`{"message" : "Server is working fine", "status": "%v", "time" : "%v"}`, 200, time.DateTime)))
		if err != nil {

		}
	})

	conn := http.Server{
		Addr:    cfg.PORT,
		Handler: server,
		ReadTimeout: cfg.READTIMEOUT,
		WriteTimeout: cfg.WRITETIMEOUT,
		IdleTimeout : cfg.IDLETIMEOUT,
	}

	fmt.Println("Welcome to the OLX API ")
	err := conn.ListenAndServe()

	if err != nil {
		fmt.Printf("%v", err.Error())
		log.Fatalf("Server Failed: %v", err.Error())
	}
}
