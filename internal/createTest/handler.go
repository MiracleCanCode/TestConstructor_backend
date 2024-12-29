package createTest

import (
	"net/http"

	"github.com/MiracleCanCode/zaperr"
	"github.com/gorilla/mux"
	"github.com/server/pkg/db"
	"github.com/server/pkg/jsonDecodeAndEncode"
	"github.com/server/pkg/middleware"
	"go.uber.org/zap"
)

type createTestHandler struct {
	logger       *zap.Logger
	db           *db.Db
	handleErrors *zaperr.Zaperr
}

func NewCreateTestHandler(logger *zap.Logger, db *db.Db, router *mux.Router, handleErrors *zaperr.Zaperr) {
	handler := &createTestHandler{
		logger:       logger,
		db:           db,
		handleErrors: handleErrors,
	}

	router.HandleFunc("/api/createAnonymusTest", handler.CreateAnonymusTest()).Methods("POST")
	router.HandleFunc("/api/createTest", middleware.IsAuthMiddleware(handler.CreateTest())).Methods("POST")
}

func (s *createTestHandler) CreateAnonymusTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		json := jsonDecodeAndEncode.NewDecodeAndEncodeJson(r, s.logger, w)

		var payload CreateAnonymusTestRequest

		err := s.handleErrors.LogError(json.DecodeAndValidationBody(&payload), "Failed to decode body", func() {
			http.Error(w, "Failed to decode body", http.StatusBadRequest)
		})
		if err != nil {
			return
		}

		testModel := MapCreateAnonymusTestRequestToModel(&payload)

		createTestService := NewCreateTestService(s.db, s.logger)
		createTestErr := createTestService.CreateTest(testModel)

		if createTestErr != nil {
			http.Error(w, createTestErr.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("Test created successfully"))
	}
}

func (s *createTestHandler) CreateTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		json := jsonDecodeAndEncode.NewDecodeAndEncodeJson(r, s.logger, w)

		var payload CreateTestRequest

		err := s.handleErrors.LogError(json.DecodeAndValidationBody(&payload), "Failed to decode body", func() {
			http.Error(w, "Failed to decode body", http.StatusBadRequest)
		})
		if err != nil {
			return
		}

		testModel := MapCreateTestRequestToModel(&payload)

		createTestService := NewCreateTestService(s.db, s.logger)
		createTestErr := createTestService.CreateTest(testModel)

		if createTestErr != nil {
			http.Error(w, createTestErr.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("Test created successfully"))
	}
}
