package http

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/server/adapters/storage/redis"
	"github.com/server/configs"
	"github.com/server/entity"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/internal/transport/http/middleware"
	"github.com/server/internal/usecases"
	cachemanager "github.com/server/pkg/cacheManager"
	"github.com/server/pkg/constants"
	errorshandler "github.com/server/pkg/errorsHandler"
	"github.com/server/pkg/json"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CacheManagerInterface interface {
	Get(key string, out interface{}) error
	Set(key string, value interface{}, ttl time.Duration)
	Delete(pattern string) error
}

type UserRepoInterfaceV2 interface {
	GetUserByLogin(login string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	UpdateUser(user *dtos.UpdateUserRequest) error
	CreateUser(user entity.User) error
	DeleteRefreshToken(login string) error
}

type UserUseCaseInterface interface {
	UpdateUserData(user dtos.UpdateUserRequest) error
	Logout(w http.ResponseWriter, r *http.Request, logger *zap.Logger) error
	FindUserByLogin(login string) (*entity.User, error)
}

type User struct {
	cfg        *configs.Config
	db         *gorm.DB
	logger     *zap.Logger
	router     *mux.Router
	repository UserRepoInterfaceV2
	cache      CacheManagerInterface
	usecase    UserUseCaseInterface
}

func NewUserHandler(logger *zap.Logger, db *gorm.DB, router *mux.Router, cfg *configs.Config) {
	repo := repository.NewUser(db, logger)
	rdb := redis.New()
	cache := cachemanager.New(rdb, logger)
	usecase := usecases.NewUser(repo, cache)
	handler := &User{
		logger:     logger,
		db:         db,
		router:     router,
		repository: repo,
		cfg:        cfg,
		cache:      cache,
		usecase:    usecase,
	}

	router.HandleFunc("/api/user/getData", middleware.IsAuth(handler.GetUserData())).Methods(http.MethodGet)
	router.HandleFunc("/api/user/update", middleware.IsAuth(handler.UpdateUser())).Methods(http.MethodPost)
	router.HandleFunc("/api/user/getByLogin", middleware.IsAuth(handler.GetUserByLogin())).Methods(http.MethodPost)
	router.HandleFunc("/api/user/logout", middleware.IsAuth(handler.Logout())).Methods(http.MethodGet)
}

func (s *User) GetUserData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		errorHandler := errorshandler.New(s.logger, w, r)
		JWT := jwt.NewJwt(s.logger)
		jsonHelper := json.New(r, s.logger, w)

		login, err := JWT.ExtractUserFromToken(r)
		if err != nil {
			s.logger.Error("GetUserData: failed extract user from token", zap.Error(err))
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		user, err := s.usecase.FindUserByLogin(login)
		if err != nil {
			s.logger.Error("GetUserData: failed find user by login", zap.Error(err))
			errorHandler.HandleError(constants.ErrGetUserData, http.StatusNotFound, err)
			return
		}

		modifiedUser := dtos.ToGetUserByLoginResponse(user)

		if err := jsonHelper.Encode(200, modifiedUser); err != nil {
			s.logger.Error("GetUserData: failed encode user data", zap.Error(err))
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *User) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.UpdateUserRequest
		errorHandler := errorshandler.New(s.logger, w, r)
		json := json.New(r, s.logger, w)

		if err := json.Decode(&payload); err != nil {
			s.logger.Error("UpdateUser: failed decode and validation request body", zap.Error(err))
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		err := s.usecase.UpdateUserData(payload)
		if err != nil {
			s.logger.Error("UpdateUser: failed update user data", zap.Error(err))
			errorHandler.HandleError(constants.ErrorUpdateUserData, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *User) GetUserByLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.GetUserByLoginRequest
		errorHandler := errorshandler.New(s.logger, w, r)
		json := json.New(r, s.logger, w)

		if err := json.Decode(&payload); err != nil {
			s.logger.Error("GetUserByLogin: failed decode and validation body", zap.Error(err))
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		user, err := s.usecase.FindUserByLogin(payload.Login)
		if err != nil {
			s.logger.Error("GetUserByLogin: failed find user by login", zap.Error(err))
			errorHandler.HandleError(constants.NotFoundUser, http.StatusNotFound, err)
			return
		}

		modifiedUser := dtos.ToGetUserByLoginResponse(user)
		if err := json.Encode(200, modifiedUser); err != nil {
			s.logger.Error("GetUserByLogin: failed encode response body", zap.Error(err))
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *User) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if err := s.usecase.Logout(w, r, s.logger); err != nil {
			s.logger.Error("Logout: failed logout", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
