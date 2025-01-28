package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/internal/usecases"
	"github.com/server/pkg/constants"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/errors"
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
	authUsecase := usecases.NewAuth(userRepo, authRepo, logger, jwtService, cfg)

	handler := &AuthHandler{
		logger:      logger,
		authUsecase: authUsecase,
		userRepo:    userRepo,
	}

	router.HandleFunc("/api/auth/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/auth/registration", handler.Registration).Methods(http.MethodPost)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.logger.Warn("Failed to close request body", zap.Error(err))
		}
	}()
	errorHandler := errors.New(h.logger, w, r)
	var payload dtos.LoginRequest
	jsonUtil := json.New(r, h.logger, w)
	if err := jsonUtil.DecodeAndValidationBody(&payload); err != nil {
		errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
		return
	}
	token, err := h.authUsecase.Login(&payload, w, r)
	if err != nil {
		errorHandler.HandleError(constants.LoginOrPasswordIncorrect, http.StatusBadGateway, err)
		return
	}

	if err := jsonUtil.Encode(http.StatusOK, token); err != nil {
		errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
	}
}

func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.logger.Warn("Failed to close request body", zap.Error(err))
		}
	}()

	errorHandler := errors.New(h.logger, w, r)

	var payload dtos.RegistrationRequest
	jsonUtil := json.New(r, h.logger, w)
	if err := jsonUtil.DecodeAndValidationBody(&payload); err != nil {
		errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
		return
	}

	result, err := h.authUsecase.Registration(&payload)
	if err != nil {
		errorHandler.HandleError(constants.ErrRegistration, http.StatusBadRequest, err)
		return
	}
	if err := jsonUtil.Encode(http.StatusOK, result); err != nil {
		errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
		return
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.logger.Warn("Failed to close request body", zap.Error(err))
		}
	}()

}
