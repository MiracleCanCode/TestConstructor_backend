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
	mapjson "github.com/server/pkg/mapJson"
	"github.com/server/pkg/middleware"
	"go.uber.org/zap"
)

type User struct {
	logger     *zap.Logger
	db         *postgresql.Db
	router     *mux.Router
	repository *repository.User
	cfg        *configs.Config
}

func NewUser(logger *zap.Logger, db *postgresql.Db, router *mux.Router, cfg *configs.Config) {
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
}

func (s *User) GetUserData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		errorHandler := errors.New(s.logger, w, r)
		JWT := jwt.NewJwt(s.logger)
		jsonHelper := json.New(r, s.logger, w)
		userUsecase := usecases.NewUser(s.repository, s.logger)

		login, err := JWT.ExtractUserFromCookie(r, "token")
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
	}
}

func (s *User) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.UpdateUserRequest
		errorHandler := errors.New(s.logger, w, r)
		json := json.New(r, s.logger, w)
		jsonResponses := mapjson.New(s.logger, w, r)

		if err := json.DecodeAndValidationBody(&payload); err != nil {
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		err := s.repository.UpdateUser(&payload)
		if err != nil {
			errorHandler.HandleError(constants.ErrorUpdateUserData, http.StatusInternalServerError, err)
			return
		}

		jsonResponses.JsonSuccess("Success update data!")
	}
}

func (s *User) GetUserByLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.GetUserByLoginRequest
		errorHandler := errors.New(s.logger, w, r)
		json := json.New(r, s.logger, w)
		userUsecase := usecases.NewUser(s.repository, s.logger)

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
	}
}
