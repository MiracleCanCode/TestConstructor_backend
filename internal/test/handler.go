package test

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/json"
	mapjson "github.com/server/pkg/mapJson"
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

	router.HandleFunc("/api/test/getById/{id}", handler.GetById()).Methods("GET")
	router.HandleFunc("/api/test/getAll", handler.GetAll()).Methods("POST")
}

func (s *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var (
			payload *GetAllTestsRequest
		)
		decoderAndEncoder := json.New(r, s.logger, w)
		jsonError := mapjson.New(s.logger, w, r)

		if err := decoderAndEncoder.Decode(&payload); err != nil {
			s.logger.Error("Failed to decode data")
			jsonError.JsonError(err.Error())

			return
		}

		getTests, count, err := s.service.GetAll(payload.Login, payload.Limit, payload.Offset)
		if err != nil {
			s.logger.Error("Failed to get tests")
			jsonError.JsonError("Failed to get tests: "+err.Error())

			return
		}
		tests := SetGetAllTests(getTests, count)

		if err := decoderAndEncoder.Encode(http.StatusOK, tests); err != nil {
			s.logger.Error("Failed to encode data")
			jsonError.JsonError(err.Error())
			return
		}
	}
}

func (s *Handler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		decoderAndEncoder := json.New(r, s.logger, w)
		jsonError := mapjson.New(s.logger, w, r)

		id := mux.Vars(r)["id"]
		parseId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			s.logger.Error("Failed parse id, error:" + err.Error())
			return
		}
		getTest, err := s.service.GetById(uint(parseId))
		if err != nil {
			s.logger.Error("Failed to get test")
			jsonError.JsonError("Failed to get test: "+err.Error())

			return
		}

		if err := decoderAndEncoder.Encode(http.StatusOK, getTest); err != nil {
			s.logger.Error("Failed to encode data")
			jsonError.JsonError( err.Error())
			return
		}
	}
}

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

		
		err = s.service.Create(testModel)

		if err != nil {
			jsonError.JsonError(err.Error())
			return
		}

		jsonError.JsonSuccess("Test created successfully")
	}
}


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

		
		err = s.service.Create(testModel)

		if err != nil {
			jsonError.JsonError(err.Error())
			return
		}

		jsonError.JsonSuccess("Test created successfully")
	}
}
