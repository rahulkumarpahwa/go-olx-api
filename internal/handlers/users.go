package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/go-olx-api/internal/dto/users"
	"github.com/rahulkumarpahwa/go-olx-api/internal/httpx"
	"github.com/rahulkumarpahwa/go-olx-api/internal/jwt"
	"github.com/rahulkumarpahwa/go-olx-api/internal/middleware"
)

func (uh *UserHandlers) Signup(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var body users.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		uh.Logger.Error("failed to decode", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.MalformedJSON)
		return
	}

	if err := body.Validate(); err != nil {
		var verr *users.ValidationError
		errors.As(err, &verr)
		uh.Logger.Error("missing or invalid fields", "request_id", requestId, "err", err)
		httpx.ValidationError(w, http.StatusBadRequest, err.Error(), httpx.ValidationFailed, verr.Field)
		return
	}

	userId, err := uh.UserServices.Signup(ctx, body, requestId)
	if err != nil {
		uh.Logger.Error("signup failed", "request_id", requestId, "err", err, "user_id", userId)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.InternalError)
		return
	}

	// setting tokens
	// access token
	accessToken, err := jwt.GenerateJWT(uh.Config, userId, jwt.Access, time.Minute*30) // half hour

	if err != nil {
		uh.Logger.Error("access_token failed", "request_id", requestId, "err", err, "user_id", userId)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.InternalError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true, // HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/*",
		MaxAge:   7 * 24 * 60 * 60,
	})

	// refresh token
	refreshToken, err := jwt.GenerateJWT(uh.Config, userId, jwt.Refresh, 7*24*time.Hour) // 7 days
	if err != nil {
		uh.Logger.Error("refresh_token failed", "request_id", requestId, "err", err, "user_id", userId)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.InternalError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true, // HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/auth/refresh",
		MaxAge:   30 * 60,
	})

	uh.Logger.Info("user signup successfully", "user_id", userId)

	httpx.Write(w, http.StatusAccepted, "user signup successfully")
}

func (uh *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var body users.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		uh.Logger.Error("failed to decode", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.MalformedJSON)
		return
	}

	if err := body.Validate(); err != nil {
		var verr *users.ValidationError
		errors.As(err, &verr)
		uh.Logger.Error("missing or invalid fields", "request_id", requestId, "err", err)
		httpx.ValidationError(w, http.StatusBadRequest, err.Error(), httpx.ValidationFailed, verr.Field)
		return
	}

	user, err := uh.UserServices.Login(ctx, body, requestId)
	if err != nil {
		uh.Logger.Error("login failed", "request_id", requestId, "err", err, "user_id", user.ID)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.InternalError)
		return
	}

	if user.ID == uuid.Nil {
		uh.Logger.Error("user does not exists", "request_id", requestId, "err", err, "user_id", user.ID)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.InternalError)
		return
	}

	// setup the token here

	uh.Logger.Info("user login successfully", "user_id", user.ID)

	httpx.Write(w, http.StatusAccepted, "user login successfully")
}
