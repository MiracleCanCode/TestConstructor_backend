package createTest

import (
	"net/http"

	"github.com/MiracleCanCode/zaperr"
	"github.com/gorilla/mux"
	"github.com/server/pkg/db"
	"github.com/server/pkg/jsonDecodeAndEncode"
	"go.uber.org/zap"
)

type createTestHandler struct {
	logger       *zap.Logger
	db           *db.Db
	repository   *CreateTestRepository
	handleErrors *zaperr.Zaperr
}

func NewCreateTestHandler(logger *zap.Logger, db *db.Db, router *mux.Router, handleErrors *zaperr.Zaperr) {
	handler := &createTestHandler{
		logger:       logger,
		db:           db,
		repository:   NewCreateTestRepository(db),
		handleErrors: handleErrors,
	}

	router.HandleFunc("/api/createTest", handler.CreateTest()).Methods("POST")
}

func (s *createTestHandler) CreateTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		json := jsonDecodeAndEncode.NewDecodeAndEncodeJson(r, s.logger, w)

		var payload *CreateTestRequest

		if err := json.DecodeAndValidationBody(payload); err != nil {
			http.Error(w, "Failed decode body", 400)
			return
		}

	}
}
