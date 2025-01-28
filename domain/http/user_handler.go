package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
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
		JWT := jwt.NewJwt(s.logger)
		jsonHelper := json.New(r, s.logger, w)

		login, err := JWT.ExtractUserFromCookie(r, "token")
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to extract login from cookie", constants.InternalServerError)
			return
		}

		user, err := s.repository.GetUserByLogin(login)
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to get user by login", constants.InternalServerError)
			return
		}

		modifiedUser := dtos.ToGetUserByLoginResponse(user)

		if err := jsonHelper.Encode(200, modifiedUser); err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to encode user data", constants.InternalServerError)
			return
		}
	}
}

func (s *User) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.UpdateUserRequest
		json := json.New(r, s.logger, w)
		jsonResponses := mapjson.New(s.logger, w, r)

		if err := json.DecodeAndValidationBody(&payload); err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to decode data", constants.InternalServerError)
			return
		}

		err := s.repository.UpdateUser(&payload)
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to update data", constants.ErrorUpdateUserData)
			return
		}

		jsonResponses.JsonSuccess("Success update data!")
	}
}

func (s *User) GetUserByLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.GetUserByLoginRequest
		json := json.New(r, s.logger, w)

		if err := json.DecodeAndValidationBody(&payload); err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to decode payload body", constants.InternalServerError)
			return
		}

		user, err := s.repository.GetUserByLogin(payload.Login)
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to get user by login", constants.NotFoundUser)
			return
		}

		if err := json.Encode(200, user); err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to encode data", constants.InternalServerError)
			return
		}
	}
}
