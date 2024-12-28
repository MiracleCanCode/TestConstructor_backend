package auth

import (
	"net/http"

	"github.com/MiracleCanCode/zaperr"
	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/pkg/db"
	"github.com/server/pkg/jsonDecodeAndEncode"
	"go.uber.org/zap"
)

type AuthHandler struct {
	log          *zap.Logger
	db           *db.Db
	cfg          *configs.Config
	service      *AuthService
	handleErrors *zaperr.Zaperr
}

func NewAuthHandler(router *mux.Router, log *zap.Logger, db *db.Db, cfg *configs.Config, handleErrors *zaperr.Zaperr) {
	handler := &AuthHandler{
		log:          log,
		db:           db,
		cfg:          cfg,
		service:      NewAuthService(db, log, cfg, handleErrors),
		handleErrors: handleErrors,
	}

	router.HandleFunc("/api/login", handler.Login()).Methods("POST")
	router.HandleFunc("/api/registration", handler.Registration()).Methods("POST")
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		json := jsonDecodeAndEncode.NewDecodeAndEncodeJson(r, h.log, w)

		var payload *LoginRequest
		if err := json.Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		user, err := h.service.Login(payload)
		if err != nil {
			http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
			h.log.Warn("Login failed", zap.Error(err))
			return
		}

		h.handleErrors.LogError(json.Encode(http.StatusOK, &user), "Failed encode data")
	}
}

func (h *AuthHandler) Registration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		json := jsonDecodeAndEncode.NewDecodeAndEncodeJson(r, h.log, w)

		var payload RegistrationRequest
		if err := json.Decode(payload); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		findUserByLogin, err := h.service.repo.GetUserByLogin(payload.Login)
		if err == nil && findUserByLogin != nil {
			http.Error(w, "Login is already taken", http.StatusBadRequest)
			h.log.Warn("User login is already taken", zap.String("login", payload.Login))
			return
		}

		user, err := h.service.Registration(&payload)
		if err != nil {
			http.Error(w, "Registration failed: "+err.Error(), http.StatusBadRequest)
			h.log.Error("Registration failed", zap.Error(err))
			return
		}

		if err := json.Encode(http.StatusOK, user); err != nil {
			h.log.Error("Failed to encode response", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
