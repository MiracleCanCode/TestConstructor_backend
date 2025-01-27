package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/pkg/cookie"
	"github.com/server/pkg/db/postgresql"
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
		cookies := cookie.New(w, r, s.logger)
		jsonHelper := json.New(r, s.logger, w)
		jsonData := mapjson.New(s.logger, w, r)

		token := cookies.Get("token")

		_, claims, err := JWT.VerifyToken(token)
		if err != nil {
			return
		}

		userLogin := claims["login"].(string)

		user, err := s.repository.GetUserByLogin(userLogin)
		if err != nil {
			s.logger.Error("Failed to get user by login", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		modifiedUser := dtos.ToGetUserByLoginResponse(user)

		if err := jsonHelper.Encode(200, modifiedUser); err != nil {
			s.logger.Error("Failed to encode user data", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			jsonData.JsonError("Failed to encode user data: " + err.Error())
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
			jsonResponses.JsonError("Failed to decode data, error:" + err.Error())
			s.logger.Error("Failed to decode data", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		err := s.repository.UpdateUser(&payload)
		if err != nil {
			jsonResponses.JsonError("Failed to update data, error:" + err.Error())
			s.logger.Error("Failed to update data", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
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
		jsonResponses := mapjson.New(s.logger, w, r)

		if err := json.DecodeAndValidationBody(&payload); err != nil {
			s.logger.Error("Failed to decode payload body", zap.Error(err))
			return
		}

		user, err := s.repository.GetUserByLogin(payload.Login)
		if err != nil {
			s.logger.Error("Failed to get user by login", zap.Error(err))
			jsonResponses.JsonError("Failed to get user by login, error" + err.Error())
			return
		}

		if err := json.Encode(200, user); err != nil {
			s.logger.Error("Failed to encode data", zap.Error(err))
			return
		}
	}
}
