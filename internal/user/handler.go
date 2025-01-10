package user

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/json"
	"github.com/server/pkg/jwt"
	mapjson "github.com/server/pkg/mapJson"
	"github.com/server/pkg/middleware"
	"go.uber.org/zap"
)

type Handler struct {
	logger     *zap.Logger
	db         *postgresql.Db
	router     *mux.Router
	repository *Repository
	cfg        *configs.Config
}

func New(logger *zap.Logger, db *postgresql.Db, router *mux.Router, cfg *configs.Config) {
	handler := &Handler{
		logger:     logger,
		db:         db,
		router:     router,
		repository: NewRepository(db, logger),
		cfg:        cfg,
	}

	router.HandleFunc("/api/user/getData", middleware.IsAuth(handler.GetData())).Methods("GET")
	router.HandleFunc("/api/user/update", middleware.IsAuth(handler.Update())).Methods("POST")
}

func (s *Handler) GetData() http.HandlerFunc {
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

		user, err := s.repository.GetByLogin(userLogin)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := jsonHelper.Encode(200, user); err != nil {
			s.logger.Error("Failed to encode user data", zap.Error(err))
			jsonData.JsonError("Failed to encode user data: " + err.Error())
			return
		}
	}
}

func (s *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload UpdateRequest
		json := json.New(r, s.logger, w)
		jsonResponses := mapjson.New(s.logger, w, r)

		if err := json.DecodeAndValidationBody(&payload); err != nil {
			jsonResponses.JsonError("Failed to decode data, error:" + err.Error())
			s.logger.Error("Failed to decode data, error:", zap.Error(err))
			return
		}

		err := s.repository.Update(&payload)
		if err != nil {
			jsonResponses.JsonError("Failed to update data, error:" + err.Error())
			s.logger.Error("Failed to update data, error:", zap.Error(err))
			return
		}

		jsonResponses.JsonSuccess("Success update data!")
	}
}
