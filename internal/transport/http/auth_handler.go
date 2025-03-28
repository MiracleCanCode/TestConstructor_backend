package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/internal/usecases"
	"github.com/server/pkg/constants"
	cookiesmanager "github.com/server/pkg/cookiesManager"
	errorshandler "github.com/server/pkg/errorsHandler"
	"github.com/server/pkg/json"
	"github.com/server/pkg/jwt"
	"gorm.io/gorm"

	"go.uber.org/zap"
)

type AuthUseCaseInterface interface {
	Login(data *dtos.LoginRequest, w http.ResponseWriter, r *http.Request) (*dtos.LoginResponse, error)
	Registration(data *dtos.RegistrationRequest) (*dtos.RegistrationResponse, error)
}

type AuthHandler struct {
	logger      *zap.Logger
	authUsecase AuthUseCaseInterface
}

func NewAuthHandler(router *mux.Router, logger *zap.Logger, db *gorm.DB, cfg *configs.Config) {
	userRepo := repository.NewUser(db, logger)
	authRepo := repository.NewAuth(db, logger)
	jwtService := jwt.NewJwt(logger)
	authUsecase := usecases.NewAuth(userRepo, authRepo, jwtService, cfg)

	handler := &AuthHandler{
		logger:      logger,
		authUsecase: authUsecase,
	}

	router.HandleFunc("/auth/login", handler.Login()).Methods(http.MethodPost)
	router.HandleFunc("/auth/registration", handler.Registration()).Methods(http.MethodPost)
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		cookie := cookiesmanager.New(r, h.logger)
		errorHandler := errorshandler.New(h.logger, w, r)
		var payload dtos.LoginRequest
		jsonUtil := json.New(r, h.logger, w)

		if err := jsonUtil.DecodeAndValidationBody(&payload); err != nil {
			h.logger.Error("Login: failed decode and validation request body", zap.Error(err))
			errorHandler.HandleError(constants.InternalServerError, http.StatusBadRequest, err)
			return
		}

		token, err := h.authUsecase.Login(&payload, w, r)
		if err != nil {
			h.logger.Error("Login: failed user login", zap.Error(err))
			errorHandler.HandleError(constants.LoginOrPasswordIncorrect, http.StatusUnauthorized, err)
			return
		}

		cookie.Set("token", token.Token, time.Minute*15, true, w)
		if err := jsonUtil.Encode(http.StatusOK, token); err != nil {
			h.logger.Error("Login: failed encode response body", zap.Error(err))
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}
	}
}

func (h *AuthHandler) Registration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		errorHandler := errorshandler.New(h.logger, w, r)
		var payload dtos.RegistrationRequest
		jsonUtil := json.New(r, h.logger, w)

		if err := jsonUtil.DecodeAndValidationBody(&payload); err != nil {
			h.logger.Error("Registration: failed decode and validation request body", zap.Error(err))
			errorHandler.HandleError(constants.InternalServerError, http.StatusBadRequest, err)
			return
		}

		result, err := h.authUsecase.Registration(&payload)
		if err != nil {
			h.logger.Error("Registration: failed registration user", zap.Error(err))
			errorHandler.HandleError(constants.ErrRegistration, http.StatusBadRequest, err)
			return
		}

		if err := jsonUtil.Encode(http.StatusCreated, result); err != nil {
			h.logger.Error("Registration: failed encode response body", zap.Error(err))
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}
	}
}
