package httpx

import (
	"encoding/json"
	"net/http"
)

/*
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
*/

// we are doing this way so that no one can pass the any random string and any var of type string as in ErrorCode can pass the values other than defined in these.
type errorCode struct {
	value string
}

var (
	InvalidId        = errorCode{value: "invalid_id"}        // 400
	NotFound         = errorCode{value: "not_found"}         // 404
	InternalError    = errorCode{value: "internal_error"}    // 500
	MalformedJSON    = errorCode{value: "malformed_json"}    // 400
	ValidationFailed = errorCode{value: "validation_failed"} // 422
	UnAuthenticated  = errorCode{value: "unauthenticated"}   // 401
	Forbidden        = errorCode{value: "forbidden"}         // 403
	Conflict         = errorCode{value: "conflict"}          // 409
	RateLimited      = errorCode{value: "rate_limited"}      // 429
)

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func Error(w http.ResponseWriter, status int, message string, code errorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(errorEnvelope{
		Error: errorPayload{
			Message: message,
			Code:    code.value,
		},
	})
}
