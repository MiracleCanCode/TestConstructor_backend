package getTest

import (
	"net/http"

	"github.com/MiracleCanCode/zaperr"
	"github.com/gorilla/mux"
	"github.com/server/pkg/db"
	"github.com/server/pkg/jsonDecodeAndEncode"
	"go.uber.org/zap"
)

type getTestHandler struct {
	logger       *zap.Logger
	db           *db.Db
	handleErrors *zaperr.Zaperr
	service      *GetTestService
}

func NewGetTestHandler(logger *zap.Logger, db *db.Db, router *mux.Router, handleErrors *zaperr.Zaperr) {
	handler := &getTestHandler{
		logger:       logger,
		db:           db,
		handleErrors: handleErrors,
		service:      NewGetTestService(db, logger, NewGetTestRepository(db)),
	}

	router.HandleFunc("/api/getTestById/{id}", handler.GetTestById()).Methods("GET")
	router.HandleFunc("/api/getAllTests", handler.GetAllTests()).Methods("POST")
}

func (s *getTestHandler) GetAllTests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var (
			payload *GetAllTestsRequest
		)
		decoderAndEncoder := jsonDecodeAndEncode.NewDecodeAndEncodeJson(r, s.logger, w)

		if err := decoderAndEncoder.Decode(&payload); err != nil {
			s.handleErrors.LogError(err, "Failed to decode data", func() {
				http.Error(w, err.Error(), http.StatusBadRequest)
			})
			return
		}

		getTests, count, err := s.service.GetAllTests(payload.Login, payload.Limit, payload.Offset)
		if err != nil {
			s.handleErrors.LogError(err, "Failed to get tests", func() {
				http.Error(w, "Failed to get tests: "+err.Error(), http.StatusInternalServerError)
			})
			return
		}
		tests := SetDataToGetAllTestsResponse(getTests, count)

		if err := decoderAndEncoder.Encode(http.StatusOK, tests); err != nil {
			s.handleErrors.LogError(err, "Failed to encode data", func() {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			})
		}
	}
}

func (s *getTestHandler) GetTestById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		decoderAndEncoder := jsonDecodeAndEncode.NewDecodeAndEncodeJson(r, s.logger, w)

		id := mux.Vars(r)["id"]
		getTest, err := s.service.GetTestById(id)
		if err != nil {
			s.handleErrors.LogError(err, "Failed to get test", func() {
				http.Error(w, "Failed to get test: "+err.Error(), http.StatusInternalServerError)
			})
			return
		}

		if err := decoderAndEncoder.Encode(http.StatusOK, getTest); err != nil {
			s.handleErrors.LogError(err, "Failed to encode data", func() {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			})
		}
	}
}
