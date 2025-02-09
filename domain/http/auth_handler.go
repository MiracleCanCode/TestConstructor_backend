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
	"github.com/server/pkg/db/postgresql"
	errorshandler "github.com/server/pkg/errorsHandler"
	"github.com/server/pkg/json"
	"github.com/server/pkg/jwt"

	"go.uber.org/zap"
)

type AuthHandler struct {
	logger      *zap.Logger
	authUsecase usecases.AuthInterface
}

func NewAuthHandler(router *mux.Router, logger *zap.Logger, db *postgresql.Db, cfg *configs.Config) {
	userRepo := repository.NewUser(db, logger)
	authRepo := repository.NewAuth(db, logger)
	jwtService := jwt.NewJwt(logger)
	authUsecase := usecases.NewAuth(userRepo, userRepo, authRepo, logger, jwtService, cfg)

	handler := &AuthHandler{
		logger:      logger,
		authUsecase: authUsecase,
	}

	router.HandleFunc("/api/auth/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/auth/registration", handler.Registration).Methods(http.MethodPost)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cookie := cookiesmanager.New(r, h.logger)
	errorHandler := errorshandler.New(h.logger, w, r)
	var payload dtos.LoginRequest
	jsonUtil := json.New(r, h.logger, w)

	if err := jsonUtil.DecodeAndValidationBody(&payload); err != nil {
		errorHandler.HandleError(constants.InternalServerError, http.StatusBadRequest, err)
		return
	}

	token, err := h.authUsecase.Login(&payload, w, r)
	if err != nil {
		errorHandler.HandleError(constants.LoginOrPasswordIncorrect, http.StatusUnauthorized, err)
		return
	}

	cookie.Set("token", token.Token, time.Minute*15, true, w)
	jsonUtil.Encode(http.StatusOK, token)
}

func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	errorHandler := errorshandler.New(h.logger, w, r)
	var payload dtos.RegistrationRequest
	jsonUtil := json.New(r, h.logger, w)

	if err := jsonUtil.DecodeAndValidationBody(&payload); err != nil {
		errorHandler.HandleError(constants.InternalServerError, http.StatusBadRequest, err)
		return
	}

	result, err := h.authUsecase.Registration(&payload)
	if err != nil {
		errorHandler.HandleError(constants.ErrRegistration, http.StatusBadRequest, err)
		return
	}

	jsonUtil.Encode(http.StatusCreated, result)
}
