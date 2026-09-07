package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	requestId = "X-Request-ID"
)
// this requestid to make every request that came from client logggable and findable.
func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestId)
		if id == "" {
			id = uuid.NewString()
		}
		w.Header().Add(requestId, id) // we will put the request id back in the header.
		ctx := context.WithValue(r.Context(), "requestCtxId", id) // we will also put in the context.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
