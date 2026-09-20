package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/go-olx-api/internal/dto/users"
	"github.com/rahulkumarpahwa/go-olx-api/internal/httpx"
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

	// setup the token here

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
