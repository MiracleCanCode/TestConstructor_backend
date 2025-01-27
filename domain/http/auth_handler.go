package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/internal/usecases"
	"github.com/server/pkg/cookie"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/json"
	"github.com/server/pkg/jwt"

	"go.uber.org/zap"
)

type AuthHandler struct {
	logger      *zap.Logger
	authUsecase usecases.AuthInterface
	userRepo    repository.UserInterface
}

func NewAuthHandler(router *mux.Router, logger *zap.Logger, db *postgresql.Db, cfg *configs.Config) {
	userRepo := repository.NewUser(db, logger)
	authRepo := repository.NewAuth(db, logger)
	jwtService := jwt.NewJwt(logger)
	authUsecase := usecases.NewAuth(userRepo, authRepo, logger, *jwtService, cfg)

	handler := &AuthHandler{
		logger:      logger,
		authUsecase: authUsecase,
		userRepo:    userRepo,
	}

	router.HandleFunc("/api/auth/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/auth/registration", handler.Registration).Methods(http.MethodPost)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.handleRequest(w, r, func(payload interface{}) (interface{}, error) {
		req := payload.(*dtos.LoginRequest)
		cookies := cookie.New(w, r, h.logger)
		token, err := h.authUsecase.Login(req)
		if err != nil {
			h.logger.Error("Failed login", zap.Error(err))
		}

		cookies.Set("token", token.Token)

		return token, nil
	}, &dtos.LoginRequest{})
}

func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	h.handleRequest(w, r, func(payload interface{}) (interface{}, error) {
		req := payload.(*dtos.RegistrationRequest)
		existingUser, _ := h.userRepo.GetUserByLogin(req.Login)
		if existingUser != nil {
			return nil, usecases.ErrLoginAlreadyTaken
		}

		return h.authUsecase.Registration(req)
	}, &dtos.RegistrationRequest{})
}

func (h *AuthHandler) handleRequest(
	w http.ResponseWriter,
	r *http.Request,
	handlerFunc func(interface{}) (interface{}, error),
	payload interface{},
) {
	defer r.Body.Close()

	jsonUtil := json.New(r, h.logger, w)
	if err := jsonUtil.DecodeAndValidationBody(payload); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	result, err := handlerFunc(payload)
	if err != nil {
		h.handleBusinessError(w, err)
		return
	}

	if err := jsonUtil.Encode(http.StatusOK, result); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to encode response", err)
	}
}

func (h *AuthHandler) handleBusinessError(w http.ResponseWriter, err error) {
	switch err {
	case usecases.ErrInvalidCredentials:
		h.respondError(w, http.StatusUnauthorized, "Invalid login or password", err)
	case usecases.ErrLoginAlreadyTaken:
		h.respondError(w, http.StatusConflict, "Login is already taken", err)
	default:
		h.respondError(w, http.StatusInternalServerError, "Internal server error", err)
	}
}

func (h *AuthHandler) respondError(w http.ResponseWriter, status int, message string, err error) {
	h.logger.Warn(message, zap.Error(err))
	http.Error(w, message, status)
}
