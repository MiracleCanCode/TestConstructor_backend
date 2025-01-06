package validateresulttest

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/server/internal/test"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/json"
	mapjson "github.com/server/pkg/mapJson"
	"go.uber.org/zap"
)

type Handler struct {
	db      *postgresql.Db
	router  *mux.Router
	service *Service
	logger  *zap.Logger
}

func New(db *postgresql.Db, router *mux.Router, logger *zap.Logger) {
	handler := &Handler{
		db:      db,
		router:  router,
		logger:  logger,
		service: NewService(db, logger, test.NewService(db, logger, test.NewRepository(db))),
	}

	handler.router.HandleFunc("/api/test/validate", handler.ValidateResult()).Methods("POST")
}


func (s *Handler) ValidateResult() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload *RequestPayload
		jsonDecodeAndEncode := json.New(r, s.logger, w)
		jsonResponse := mapjson.New(s.logger, w, r)
		if err := jsonDecodeAndEncode.Decode(&payload); err != nil {
			s.logger.Error("Failed to decode body: " + err.Error())
			jsonResponse.JsonError("Invalid request payload")
			return
		}
		result, err := s.service.Validate(payload.Test)
		if err != nil {
			s.logger.Error("Validation failed: " + err.Error())
			jsonResponse.JsonError("Failed to validate test")
			return
		}

		jsonResponse.JsonSuccess("percents: " + strconv.FormatFloat(*result, 'f', 2, 64))

	}
}

