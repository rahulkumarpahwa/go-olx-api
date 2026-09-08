package httpx

import (
	"encoding/json"
	"net/http"
)

type ErrorCode string

const (
	InvalidId        ErrorCode = "invalid_id"        // 400
	NotFound         ErrorCode = "not_found"         // 404
	InternalError    ErrorCode = "internal_error"    // 500
	MalformedJSON    ErrorCode = "malformed_json"    // 400
	ValidationFailed ErrorCode = "validation_failed" // 422
	UnAuthenticated  ErrorCode = "unauthenticated"   // 401
	Forbidden        ErrorCode = "forbidden"         // 403
	Conflict         ErrorCode = "conflict"          // 409
	RateLimited      ErrorCode = "rate_limited"      // 429
)

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

func Error(w http.ResponseWriter, status int, message string, code ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(errorEnvelope{
		Error: errorPayload{
			Message: message,
			Code:    code,
		},
	})
}
