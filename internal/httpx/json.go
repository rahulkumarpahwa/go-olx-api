package httpx

import (
	"encoding/json"
	"net/http"
)

type dataEnvelope struct {
	Success bool        `json:"success"`
	Data    any `json:"data"`
}

func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(dataEnvelope{
		Success: true,
		Data: data,
	})
}
