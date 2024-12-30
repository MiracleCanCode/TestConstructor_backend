package validateresulttest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/internal/getTest"
	"github.com/server/pkg/db"
	"github.com/server/pkg/jsonDecodeAndEncode"
	"go.uber.org/zap"
)

type ValidateResultTestHandler struct {
	db      *db.Db
	router  *mux.Router
	service *ValidateResultTestService
	logger  *zap.Logger
}

func NewValidateTestHandler(db *db.Db, router *mux.Router, logger *zap.Logger) {
	handler := ValidateResultTestHandler{
		db:      db,
		router:  router,
		logger:  logger,
		service: NewValidateResultTestService(db, logger, getTest.NewGetTestService(db, logger, getTest.NewGetTestRepository(db))),
	}

	handler.router.HandleFunc("/api/validate", handler.ValidateResult()).Methods("POST")
}

func (s *ValidateResultTestHandler) ValidateResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload ValidateResultTestRequest
		jsonDecodeAndEncode := jsonDecodeAndEncode.NewDecodeAndEncodeJson(r, s.logger, w)

		if err := jsonDecodeAndEncode.Decode(&payload); err != nil {
			s.logger.Error("Failed to decode body: " + err.Error())
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		result, err := s.service.Validate(payload.Test)
		if err != nil {
			s.logger.Error("Validation failed: " + err.Error())
			http.Error(w, "Failed to validate test", http.StatusInternalServerError)
			return
		}

		if err := jsonDecodeAndEncode.Encode(http.StatusOK, result); err != nil {
			s.logger.Error("Failed to encode response: " + err.Error())
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
