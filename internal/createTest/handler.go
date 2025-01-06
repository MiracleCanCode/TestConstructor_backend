package createTest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/json"
	mapjson "github.com/server/pkg/mapJson"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	db     *postgresql.Db
}

// NewCreateTestHandler создает новые маршруты для создания тестов
// @Summary Initialize create test routes
// @Description Set up routes for creating anonymous and authenticated tests
// @Tags create-test
// @Param logger path string true "Logger"
// @Param db path string true "Database"
// @Param router path string true "Router"
// @Param handleErrors path string true "Error handler"
// @Success 200 {string} string "Routes initialized successfully"
func New(logger *zap.Logger, db *postgresql.Db, router *mux.Router) {
	handler := &Handler{
		logger: logger,
		db:     db,
	}

	router.HandleFunc("/api/createAnonymusTest", handler.AnonymousTest()).Methods("POST")
	router.HandleFunc("/api/createTest", handler.StandardTest()).Methods("POST")
}

// CreateAnonymusTest - обработчик для создания анонимного теста
// @Summary Create an anonymous test
// @Description Creates a new anonymous test using provided data
// @Tags create-test
// @Accept json
// @Produce json
// @Param test body CreateAnonymusTestRequest true "Anonymus Test Data"
// @Success 201 {string} string "Test created successfully"
// @Failure 400 {object} ErrorResponse "Failed to decode body"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/createAnonymusTest [post]
func (s *Handler) AnonymousTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload CreateAnonymusTestRequest
		json := json.New(r, s.logger, w)
		jsonError := mapjson.New(s.logger, w, r)

		err := json.DecodeAndValidationBody(&payload)
		if err != nil {
			s.logger.Error("Failed to decode body")
			jsonError.JsonError("Failed to decode body")
			return
		}

		testModel := MapCreateAnonymusTestRequestToModel(&payload)

		createTestService := NewService(s.db, s.logger)
		err = createTestService.Create(testModel)

		if err != nil {
			jsonError.JsonError(err.Error())
			return
		}

		jsonError.JsonSuccess("Test created successfully")
	}
}

// CreateTest - обработчик для создания теста
// @Summary Create a new test
// @Description Creates a new authenticated test using provided data
// @Tags create-test
// @Accept json
// @Produce json
// @Param test body CreateTestRequest true "Test Data"
// @Success 201 {string} string "Test created successfully"
// @Failure 400 {object} ErrorResponse "Failed to decode body"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/createTest [post]
func (s *Handler) StandardTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload CreateTestRequest
		json := json.New(r, s.logger, w)
		jsonError := mapjson.New(s.logger, w, r)

		err := json.DecodeAndValidationBody(&payload)
		if err != nil {
			s.logger.Error("Failed to decode body")
			jsonError.JsonError("Failed to decode body:" + err.Error())
			return
		}

		testModel := MapCreateTestRequestToModel(&payload)

		createTestService := NewService(s.db, s.logger)
		err = createTestService.Create(testModel)

		if err != nil {
			jsonError.JsonError(err.Error())
			return
		}

		jsonError.JsonSuccess("Test created successfully")
	}
}
