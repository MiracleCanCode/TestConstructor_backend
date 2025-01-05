package validateresulttest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/internal/getTest"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/json"
	"go.uber.org/zap"
)

type Handler struct {
	db      *postgresql.Db
	router  *mux.Router
	service *Service
	logger  *zap.Logger
}

// NewValidateTestHandler создаёт новый обработчик для валидации теста
// @Summary Validate a test result
// @Description Validates the test result based on the input data
// @Tags validation
// @Accept json
// @Produce json
// @Param validateRequest body ValidateResultTestRequest true "Validation request payload"
// @Success 200 {object} ValidationResult "Validation result data"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 500 {object} ErrorResponse "Failed to validate test"
// @Router /api/validate [post]
func New(db *postgresql.Db, router *mux.Router, logger *zap.Logger) {
	handler := &Handler{
		db:      db,
		router:  router,
		logger:  logger,
		service: NewService(db, logger, getTest.NewService(db, logger, getTest.NewRepository(db))),
	}

	handler.router.HandleFunc("/api/validate", handler.ValidateResult()).Methods("POST")
}

// ValidateResult - обработчик для валидации результата теста
// @Summary Validate a test result
// @Description Validates the test result based on the provided test data
// @Tags validation
// @Accept json
// @Produce json
// @Param test body ValidateResultTestRequest true "Test result data to validate"
// @Success 200 {object} ValidationResult "Validation result"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/validate [post]
func (s *Handler) ValidateResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload *RequestPayload
		jsonDecodeAndEncode := json.New(r, s.logger, w)
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

