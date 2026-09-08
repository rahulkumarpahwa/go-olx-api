package httpx

import (
	"encoding/json"
	"net/http"
)


type ErrorCode string

const (
	NoRows            ErrorCode  = "no_rows"
	InternalServerError ErrorCode = "internal_server_error"
	SomethingWentWrong  ErrorCode = "some_went_wrong"
	BadRequest ErrorCode = "bad_request"
)


type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Message string `json:"message"`
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
