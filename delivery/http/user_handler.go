package http

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/dtos"
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/internal/utils/json"
	"github.com/server/internal/utils/jwt"
	mapjson "github.com/server/internal/utils/mapJson"
	"github.com/server/internal/utils/middleware"
	"github.com/server/models"
	"github.com/server/repository"
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
}

func (s *User) GetUserData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		jsonHelper := json.New(r, s.logger, w)
		jsonData := mapjson.New(s.logger, w, r)

		authHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, claims, err := jwt.NewJwt("SUPERSECRETKEYFORBESTAPPINTHEWORLD").VerifyToken(tokenString)
		if err != nil {
			return
		}

		userLogin := claims["login"].(string)

		userChan := make(chan *models.User, 1)
		errChan := make(chan error, 1)

		go func() {
			user, err := s.repository.GetUserByLogin(userLogin)
			if err != nil {
				errChan <- err
				return
			}
			userChan <- user
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
