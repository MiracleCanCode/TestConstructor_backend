package getTest

import (
	"net/http"

	"github.com/MiracleCanCode/zaperr"
	"github.com/gorilla/mux"
	"github.com/server/pkg/db"
	"github.com/server/pkg/middleware"
	"go.uber.org/zap"
)

type getTestHandler struct {
	logger       *zap.Logger
	db           *db.Db
	handleErrors *zaperr.Zaperr
}

func NewGetTestHandler(logger *zap.Logger, db *db.Db, router *mux.Router, handleErrors *zaperr.Zaperr) {
	handler := &getTestHandler{
		logger:       logger,
		db:           db,
		handleErrors: handleErrors,
	}

	router.HandleFunc("/api/getTestById", middleware.IsAuthMiddleware(handler.GetTestById())).Methods("POST")
	router.HandleFunc("/api/getAllTests", middleware.IsAuthMiddleware(handler.GetAllTests())).Methods("POST")
}

func (s *getTestHandler) GetAllTests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
func (s *getTestHandler) GetTestById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
