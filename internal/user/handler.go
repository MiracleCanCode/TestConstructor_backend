package user

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
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
	repository *Repository
}

func New(logger *zap.Logger, db *postgresql.Db, router *mux.Router) {
	handler := &Handler{
		logger: logger,
		db: db,
		router: router,
		repository: NewRepository(db, logger),
	}

	router.HandleFunc("/api/user/getData", handler.GetData()).Methods("GET")
	router.HandleFunc("/api/user/update", handler.Update() ).Methods("POST")
}

func (s *Handler) GetData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		json := json.New(r, s.logger, w)
		jsonData := mapjson.New(s.logger, w, r)
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := "supersecretkeybysuperuser"
		

		_, claims, err := jwt.NewJwt(secret).VerifyToken(tokenString)
		if err != nil {
			jsonData.JsonError("Failed decode data")
			s.logger.Error("Failed decode data")
			return
		}


		userLogin := claims["login"].(string)
		userData,_ := s.repository.GetByLogin(userLogin)
		if err :=json.Encode(200, userData); err != nil {
			jsonData.JsonError("Failed encode data")
			s.logger.Error("Failed encode data")
			return
		}
	
	}
}

func (s *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

	}
}