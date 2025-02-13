package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/delivery/http/middleware"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/internal/usecases"
	cachemanager "github.com/server/pkg/cacheManager"
	"github.com/server/pkg/constants"
	"github.com/server/pkg/db/redis"
	errorshandler "github.com/server/pkg/errorsHandler"
	"github.com/server/pkg/json"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type User struct {
	cfg        *configs.Config
	db         *gorm.DB
	logger     *zap.Logger
	router     *mux.Router
	repository *repository.User
}

func NewUserHandler(logger *zap.Logger, db *gorm.DB, router *mux.Router, cfg *configs.Config) {
	handler := &User{
		logger:     logger,
		db:         db,
		router:     router,
		repository: repository.NewUser(db, logger),
		cfg:        cfg,
	}

	router.HandleFunc("/api/user/getData", middleware.IsAuth(handler.GetUserData())).Methods(http.MethodGet)
	router.HandleFunc("/api/user/update", middleware.IsAuth(handler.UpdateUser())).Methods(http.MethodPost)
	router.HandleFunc("/api/user/getByLogin", middleware.IsAuth(handler.GetUserByLogin())).Methods(http.MethodPost)
	router.HandleFunc("/api/user/logout", middleware.IsAuth(handler.Logout())).Methods(http.MethodGet)
}

func (s *User) GetUserData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				s.logger.Warn("Failed to close request body", zap.Error(err))
			}
		}()

		errorHandler := errorshandler.New(s.logger, w, r)
		JWT := jwt.NewJwt(s.logger)
		jsonHelper := json.New(r, s.logger, w)
		rdb := redis.New()
		cache := cachemanager.New(rdb, s.logger)
		userUsecase := usecases.NewUser(s.repository, s.logger, cache)

		login, err := JWT.ExtractUserFromToken(r)
		if err != nil {
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		user, err := userUsecase.FindUserByLogin(login)
		if err != nil {
			errorHandler.HandleError(constants.ErrGetUserData, http.StatusNotFound, err)
			return
		}

		modifiedUser := dtos.ToGetUserByLoginResponse(user)

		if err := jsonHelper.Encode(200, modifiedUser); err != nil {
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *User) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				s.logger.Warn("Failed to close request body", zap.Error(err))
			}
		}()

		var payload dtos.UpdateUserRequest
		errorHandler := errorshandler.New(s.logger, w, r)
		json := json.New(r, s.logger, w)
		rdb := redis.New()
		cache := cachemanager.New(rdb, s.logger)
		userUsecase := usecases.NewUser(s.repository, s.logger, cache)

		if err := json.DecodeAndValidationBody(&payload); err != nil {
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		err := userUsecase.UpdateUserData(payload)
		if err != nil {
			errorHandler.HandleError(constants.ErrorUpdateUserData, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *User) GetUserByLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				s.logger.Warn("Failed to close request body", zap.Error(err))
			}
		}()

		var payload dtos.GetUserByLoginRequest
		errorHandler := errorshandler.New(s.logger, w, r)
		json := json.New(r, s.logger, w)
		rdb := redis.New()
		cache := cachemanager.New(rdb, s.logger)
		userUsecase := usecases.NewUser(s.repository, s.logger, cache)

		if err := json.DecodeAndValidationBody(&payload); err != nil {
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		user, err := userUsecase.FindUserByLogin(payload.Login)
		if err != nil {
			errorHandler.HandleError(constants.NotFoundUser, http.StatusNotFound, err)
			return
		}

		modifiedUser := dtos.ToGetUserByLoginResponse(user)
		if err := json.Encode(200, modifiedUser); err != nil {
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *User) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				s.logger.Warn("Failed to close request body", zap.Error(err))
			}
		}()
		rdb := redis.New()
		cache := cachemanager.New(rdb, s.logger)
		userUsecase := usecases.NewUser(s.repository, s.logger, cache)

		if err := userUsecase.Logout(w, r); err != nil {
			s.logger.Error("Logout", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
