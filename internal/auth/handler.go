package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/configs"
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
}

// NewAuthHandler создаёт новый обработчик аутентификации
// @Summary Initialize auth routes
// @Description Set up routes for login and registration
// @Tags auth
// @Param router path string true "Router"
// @Param log path string true "Logger"
// @Param db path string true "Database"
// @Param cfg path string true "Configuration"
// @Param handleErrors path string true "Error handler"
// @Success 200 {object} string "Routes initialized successfully"
func New(router *mux.Router, log *zap.Logger, db *postgresql.Db, cfg *configs.Config) {
	handler := &Handler{
		log:     log,
		db:      db,
		cfg:     cfg,
		service: NewService(db, log, cfg),
	}

	router.HandleFunc("/api/login", handler.Login()).Methods("POST")
	router.HandleFunc("/api/registration", handler.Registration()).Methods("POST")
}

// Login - обработчик для аутентификации пользователя
// @Summary Login to the system
// @Description Authenticates the user with login and password
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Login Request"
// @Success 200 {object} LoginResponse "User data"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Router /api/login [post]
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

// Registration - обработчик для регистрации нового пользователя
// @Summary Register a new user
// @Description Registers a new user with a unique login and password
// @Tags auth
// @Accept json
// @Produce json
// @Param registrationRequest body RegistrationRequest true "Registration Request"
// @Success 200 {object} RegistrationResponse "User data"
// @Failure 400 {object} ErrorResponse "Invalid request payload or login already taken"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/registration [post]
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

		findUserByLogin, err := h.service.repo.GetUserByLogin(payload.Login)
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

