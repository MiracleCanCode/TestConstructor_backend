package http

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
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

	router.HandleFunc("/api/user/getData", middleware.IsAuth(handler.GetUserData())).Methods("GET")
	router.HandleFunc("/api/user/update", middleware.IsAuth(handler.UpdateUser())).Methods("POST")
	router.HandleFunc("/api/user/getByLogin", middleware.IsAuth(handler.GetUserByLogin())).Methods("POST")
}

func (s *User) GetUserData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		jsonHelper := json.New(r, s.logger, w)
		jsonData := mapjson.New(s.logger, w, r)

		authHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, claims, err := jwt.NewJwt(s.logger).VerifyToken(tokenString)
		if err != nil {
			return
		}

		userLogin := claims["login"].(string)

		userChan := make(chan *dtos.GetUserByLoginResponse, 1)
		errChan := make(chan error, 1)

		go func() {
			user, err := s.repository.GetUserByLogin(userLogin)
			modifiedUser := dtos.ToGetUserByLoginResponse(user)
			if err != nil {
				errChan <- err
				return
			}
			userChan <- modifiedUser
		}()

		select {
		case user := <-userChan:
			if err := jsonHelper.Encode(200, user); err != nil {
				s.logger.Error("Failed to encode user data", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
				jsonData.JsonError("Failed to encode user data: " + err.Error())
				return
			}
		case err := <-errChan:
			s.logger.Error("Failed to get user by login", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
