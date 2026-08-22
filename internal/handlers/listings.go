package handlers

import (
	"database/sql"
	"net/http"
)

func Listings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Listings"))
	}
}
