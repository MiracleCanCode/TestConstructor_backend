package user

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/server/internal/auth"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/json"
	"github.com/server/pkg/jwt"
	mapjson "github.com/server/pkg/mapJson"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	db *postgresql.Db
	router *mux.Router
}

func New(logger *zap.Logger, db *postgresql.Db, router *mux.Router) {
	handler := &Handler{
		logger: logger,
		db: db,
		router: router,
	}

	router.HandleFunc("/api/getData", handler.GetData()).Methods("GET")
}

func (s *Handler) GetData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json := json.New(r, s.logger, w)
		jsonData := mapjson.New(s.logger, w, r)
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := "supersecretkeybysuperuser"
		user := auth.NewRepository(s.db, s.logger)

		_, claims, err := jwt.NewJwt(secret).VerifyToken(tokenString)
		if err != nil {
			jsonData.JsonError("Не получилось декодировать токен")
			s.logger.Error("Не получилось декодировать токен")
			return
		}


		userLogin := claims["login"].(string)
		userData,_ := user.GetUserByLogin(userLogin)
		if err :=json.Encode(200, userData); err != nil {
			jsonData.JsonError("Не получилось отправить данные")
			s.logger.Error("Не получилось отправить данные")
			return
		}
	
	}
}