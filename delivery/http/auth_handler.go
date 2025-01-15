package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/dtos"
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/internal/utils/json"
	mapjson "github.com/server/internal/utils/mapJson"
	"github.com/server/repository"
	"github.com/server/usecases"
	"go.uber.org/zap"
)

type Auth struct {
	log      *zap.Logger
	db       *postgresql.Db
	cfg      *configs.Config
	service  *usecases.Auth
	userRepo *repository.User
}

func NewAuth(router *mux.Router, log *zap.Logger, db *postgresql.Db, cfg *configs.Config) {
	service := usecases.NewAuth(db, log, cfg, repository.NewUser(db, log))
	handler := &Auth{
		log:      log,
		db:       db,
		cfg:      cfg,
		service:  service,
		userRepo: repository.NewUser(db, log),
	}

	router.HandleFunc("/api/auth/login", handler.Login()).Methods("POST")
	router.HandleFunc("/api/auth/registration", handler.Registration()).Methods("POST")
}

func (h *Auth) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		json := json.New(r, h.log, w)
		message := mapjson.New(h.log, w, r)

		var payload dtos.LoginRequest
		if err := json.DecodeAndValidationBody(&payload); err != nil {
			message.JsonError("Invalid request payload")
			h.log.Warn("Invalid login request", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		user, err := h.service.Login(&payload)
		if err != nil {
			message.JsonError("Login failed: " + err.Error())
			h.log.Warn("Login failed", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		if err := json.Encode(http.StatusOK, &user); err != nil {
			h.log.Error("Failed to encode login response", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			message.JsonError("Internal server error")
			return
		}
	}
}

func (h *Auth) Registration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		json := json.New(r, h.log, w)
		message := mapjson.New(h.log, w, r)

		var payload dtos.RegistrationRequest
		if err := json.DecodeAndValidationBody(&payload); err != nil {
			message.JsonError("Invalid request payload")
			h.log.Warn("Invalid registration request", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		findUserByLogin, err := h.userRepo.GetUserByLogin(payload.Login)
		if err == nil && findUserByLogin != nil {
			message.JsonError("Login is already taken")
			h.log.Warn("User login is already taken", zap.String("login", payload.Login), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		user, err := h.service.Registration(&payload)
		if err != nil {
			message.JsonError("Registration failed: " + err.Error())
			h.log.Error("Registration failed", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		if err := json.Encode(http.StatusOK, user); err != nil {
			message.JsonError("Failed to encode registration response, error: " + err.Error())
			h.log.Error("Failed encode data", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}
	}
}
