package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/server/delivery/http/middleware"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/internal/usecases"
	"github.com/server/pkg/constants"
	errorshandler "github.com/server/pkg/errorsHandler"
	"github.com/server/pkg/json"
	mapjson "github.com/server/pkg/mapJson"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ValidateResult struct {
	db      *gorm.DB
	router  *mux.Router
	service usecases.TestValidatorInterface
	logger  *zap.Logger
}

func NewValidateResultHandler(db *gorm.DB, router *mux.Router, logger *zap.Logger) {
	testManagerRepo := repository.NewTestManager(db)
	handler := &ValidateResult{
		db:      db,
		router:  router,
		logger:  logger,
		service: usecases.NewTestValidator(testManagerRepo, logger),
	}

	handler.router.HandleFunc("/api/test/validate", middleware.IsAuth(handler.ValidateResult())).Methods(http.MethodPost)
}

func (s *ValidateResult) ValidateResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				s.logger.Warn("Failed to close request body", zap.Error(err))
			}
		}()

		var payload dtos.ValidateResultRequestPayload
		errorHandler := errorshandler.New(s.logger, w, r)
		jsonDecodeAndEncode := json.New(r, s.logger, w)
		jsonResponse := mapjson.New(s.logger, w, r)

		if err := jsonDecodeAndEncode.Decode(&payload); err != nil {
			errorHandler.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		result, err := s.service.Validate(payload.Test)
		if err != nil {
			errorHandler.HandleError(constants.ErrTestValidation, http.StatusBadRequest, err)
			return
		}

		jsonResponse.JsonSuccess(strconv.FormatFloat(*result, 'f', 2, 64))
	}
}
