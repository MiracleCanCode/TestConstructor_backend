package validateresulttest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/internal/getTest"
	"github.com/server/pkg/db"
	"go.uber.org/zap"
)

type ValidateResultTestHandler struct {
	db      *db.Db
	router  *mux.Router
	service *ValidateResultTestService
}

func NewValidateTestHandler(db *db.Db, router *mux.Router, logger *zap.Logger) {
	handler := ValidateResultTestHandler{
		db:      db,
		router:  router,
		service: NewValidateResultTestService(db, logger, getTest.NewGetTestService(db, logger, getTest.NewGetTestRepository(db))),
	}

	handler.router.HandleFunc("/api/validate", handler.ValidateResult()).Methods("POST")
}

func (s *ValidateResultTestHandler) ValidateResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
	}
}
