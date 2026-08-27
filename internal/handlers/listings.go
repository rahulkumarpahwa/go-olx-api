package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rahulkumarpahwa/go-olx-api/internal/types"
)

func Listings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const query = "SELECT id, title, description, price, status, city, user_id, category_id, created_at FROM listings"

		rows, err := db.QueryContext(r.Context(), query)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "db.rows: No Rows", http.StatusNoContent)
				return
			}
			http.Error(w, "db.rows: "+err.Error(), http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var listings []types.Listings

		for rows.Next() {
			var l types.Listings
			err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.Status, &l.City, &l.UserID, &l.CategoryID, &l.CreatedAt)
			if err != nil {
				http.Error(w, "rows.scan: "+err.Error(), http.StatusInternalServerError)
				return
			}
			listings = append(listings, l)
		}

		err = rows.Err()
		if err != nil {
			http.Error(w, "rows.Err(): "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string][]types.Listings{"listings": listings})
	}
}
