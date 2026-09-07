package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey int

// we are making this as enum, the next values will be increamental automatically.
// we are making the key private to package as well.
const (
	requestIDKey ctxKey = iota // we put the key in small caps so that it can't be exported and remains not be to creatable in some other package.
)

//The provided key must be comparable and should not be of type string or any other built-in type to avoid collisions between packages using context. Users of WithValue should define their own types for keys. To avoid allocating when assigning to an interface{}, context keys often have concrete type struct{}. Alternatively, exported context key variables' static type should be a pointer or interface.

const (
	requestId = "X-Request-ID"
)

// this requestid to make every request that came from client logggable and findable, in case of any error shown on the system.
func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestId)
		if id == "" {
			id = uuid.NewString()
		}
		w.Header().Add(requestId, id)                           // we will put the request id back in the header, so as send back to the client.
		ctx := context.WithValue(r.Context(), requestIDKey, id) // we will also put in the context.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// this function is to extract the requestId from the request's context using the key and we want to keep key remain in the same package, so as to avoid collision.
func RequestIDFromContext(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)
}
