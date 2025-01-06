package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/internal/user"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/json"
	mapjson "github.com/server/pkg/mapJson"
	"go.uber.org/zap"
)

type Handler struct {
	log     *zap.Logger
	db      *postgresql.Db
	cfg     *configs.Config
	service *Service
	userRepo *user.Repository
}


func New(router *mux.Router, log *zap.Logger, db *postgresql.Db, cfg *configs.Config) {
	service := NewService(db, log, cfg, user.NewRepository(db, log))
	handler := &Handler{
		log:     log,
		db:      db,
		cfg:     cfg,
		service: service,
		userRepo: user.NewRepository(db, log),
	}

	router.HandleFunc("/api/auth/login", handler.Login()).Methods("POST")
	router.HandleFunc("/api/auth/registration", handler.Registration()).Methods("POST")
}


func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		json := json.New(r, h.log, w)
		message := mapjson.New(h.log, w, r)

		var payload LoginRequest
		if err := json.Decode(&payload); err != nil {
			message.JsonError("Invalid request payload")
			h.log.Warn("Invalid login request", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		user, err := h.service.Login(&payload)
		if err != nil {
			message.JsonError("Login failed: "+err.Error())
			h.log.Warn("Login failed", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		if err := json.Encode(http.StatusOK, &user); err != nil {
			h.log.Error("Failed to encode login response", zap.Error(err))
			message.JsonError("Internal server error")
			return
		}
	}
}


func (h *Handler) Registration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		json := json.New(r, h.log, w)
		message := mapjson.New(h.log, w, r)

		var payload RegistrationRequest
		if err := json.Decode(&payload); err != nil {
			message.JsonError("Invalid request payload")
			h.log.Warn("Invalid registration request", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		findUserByLogin, err := h.userRepo.GetByLogin(payload.Login)
		if err == nil && findUserByLogin != nil {
			message.JsonError("Login is already taken")
			h.log.Warn("User login is already taken", zap.String("login", payload.Login), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		user, err := h.service.Registration(&payload)
		if err != nil {
			message.JsonError("Registration failed: "+err.Error())
			h.log.Error("Registration failed", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		if err := json.Encode(http.StatusOK, user); err != nil {
			message.JsonError("Failed to encode registration response, error: " + err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

