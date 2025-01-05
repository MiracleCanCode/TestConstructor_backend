package getTest

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/json"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.Logger
	db      *postgresql.Db
	service *Service
}

func New(logger *zap.Logger, db *postgresql.Db, router *mux.Router) {
	handler := &Handler{
		logger:  logger,
		db:      db,
		service: NewService(db, logger, NewRepository(db)),
	}

	router.HandleFunc("/api/getTestById/{id}", handler.GetById()).Methods("GET")
	router.HandleFunc("/api/getAllTests", handler.GetAll()).Methods("POST")
}

func (s *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var (
			payload *GetAllTestsRequest
		)
		decoderAndEncoder := json.New(r, s.logger, w)

		if err := decoderAndEncoder.Decode(&payload); err != nil {
			s.logger.Error("Failed to decode data")
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		getTests, count, err := s.service.GetAll(payload.Login, payload.Limit, payload.Offset)
		if err != nil {
			s.logger.Error("Failed to get tests")
			http.Error(w, "Failed to get tests: "+err.Error(), http.StatusInternalServerError)

			return
		}
		tests := SetGetAllTests(getTests, count)

		if err := decoderAndEncoder.Encode(http.StatusOK, tests); err != nil {
			s.logger.Error("Failed to encode data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Handler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		decoderAndEncoder := json.New(r, s.logger, w)

		id := mux.Vars(r)["id"]
		parseId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			s.logger.Error("Failed parse id, error:" + err.Error())
			return
		}
		getTest, err := s.service.GetById(uint(parseId))
		if err != nil {
			s.logger.Error("Failed to get test")
			http.Error(w, "Failed to get test: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := decoderAndEncoder.Encode(http.StatusOK, getTest); err != nil {
			s.logger.Error("Failed to encode data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
