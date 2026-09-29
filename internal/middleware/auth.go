package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/rahulkumarpahwa/go-olx-api/internal/config"
	"github.com/rahulkumarpahwa/go-olx-api/internal/httpx"
	"github.com/rahulkumarpahwa/go-olx-api/internal/jwt"
)

type Middleware struct {
	Config *config.Config
	Logger *slog.Logger
}

func NewMiddleware(config *config.Config, logger *slog.Logger) *Middleware {
	return &Middleware{
		Config: config,
		Logger: logger,
	}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("access_token")
		if err != nil {
			httpx.Error(w, http.StatusInternalServerError, "invalid credentials", httpx.UnAuthenticated)
			return
		}

		claims, err := jwt.VerifyToken(m.Config, cookie.Value)

		if err != nil {
			m.Logger.Error("not a verified token", "err", err)
			httpx.Error(w, http.StatusInternalServerError, "invalid credentials", httpx.UnAuthenticated)
			return
		}

		if claims.TokenType != jwt.Access {
			m.Logger.Error("not a valid token type", )
			httpx.Error(w, http.StatusInternalServerError, "invalid credentials", httpx.UnAuthenticated)
			return
		}

		// user authenticated
		ctx := context.WithValue(
			r.Context(),
			"userID",
			claims.UserID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
