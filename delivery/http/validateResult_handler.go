package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/server/dtos"
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/internal/utils/json"
	mapjson "github.com/server/internal/utils/mapJson"
	"github.com/server/internal/utils/middleware"
	"github.com/server/repository"
	"github.com/server/usecases"
	"go.uber.org/zap"
)

type ValidateResult struct {
	db      *postgresql.Db
	router  *mux.Router
	service usecases.Validator
	logger  *zap.Logger
}

func NewValidateResult(db *postgresql.Db, router *mux.Router, logger *zap.Logger) {
	testManagerRepo := repository.NewTestManager(db)
	handler := &ValidateResult{
		db:      db,
		router:  router,
		logger:  logger,
		service: usecases.NewTestValidator(testManagerRepo, logger),
	}

	handler.router.HandleFunc("/api/test/validate", middleware.IsAuth(handler.ValidateResult())).Methods("POST")
}

func (s *ValidateResult) ValidateResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.ValidateResultRequestPayload
		jsonDecodeAndEncode := json.New(r, s.logger, w)
		jsonResponse := mapjson.New(s.logger, w, r)
		if err := jsonDecodeAndEncode.Decode(&payload); err != nil {
			s.logger.Error("Failed to decode body", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			jsonResponse.JsonError("Invalid request payload")
			return
		}
		result, err := s.service.Validate(payload.Test)
		if err != nil {
			s.logger.Error("Validation failed", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		jsonResponse.JsonSuccess(strconv.FormatFloat(*result, 'f', 2, 64))
	}
}
