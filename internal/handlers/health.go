package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message" : "Server is working fine", "status": "%v", "time" : "%v"}`, 200, time.Now().Local())))
}
