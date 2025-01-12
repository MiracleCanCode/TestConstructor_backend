package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/server/dtos"
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/internal/utils/json"
	mapjson "github.com/server/internal/utils/mapJson"
	"github.com/server/usecases"
	"go.uber.org/zap"
)

type TestManager struct {
	logger  *zap.Logger
	db      *postgresql.Db
	service *usecases.TestManager
}

func NewTestManager(logger *zap.Logger, db *postgresql.Db, router *mux.Router) {
	handler := &TestManager{
		logger:  logger,
		db:      db,
		service: usecases.NewTestManager(db, logger),
	}

	router.HandleFunc("/api/test/getById/{id}", handler.GetTestById()).Methods("GET")
	router.HandleFunc("/api/test/getAll", handler.GetAll()).Methods("POST")
	router.HandleFunc("/api/test/create", handler.CreateTest()).Methods("POST")
}

func (s *TestManager) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var (
			payload *dtos.GetAllTestsRequest
		)
		decoderAndEncoder := json.New(r, s.logger, w)
		jsonError := mapjson.New(s.logger, w, r)

		if err := decoderAndEncoder.Decode(&payload); err != nil {
			s.logger.Error("Failed to decode data")
			jsonError.JsonError(err.Error())

			return
		}

		getTests, count, err := s.service.GetAllTests(payload.Login, payload.Limit, payload.Offset)
		if err != nil {
			s.logger.Error("Failed to get tests")
			jsonError.JsonError("Failed to get tests: " + err.Error())

			return
		}
		tests := dtos.SetGetAllTests(getTests, count)

		if err := decoderAndEncoder.Encode(http.StatusOK, tests); err != nil {
			s.logger.Error("Failed to encode data")
			jsonError.JsonError(err.Error())
			return
		}
	}
}

func (s *TestManager) GetTestById() http.HandlerFunc {
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
		getTest, err := s.service.GetTestById(uint(parseId))
		if err != nil {
			s.logger.Error("Failed to get test")
			jsonError.JsonError("Failed to get test: " + err.Error())

			return
		}

		if err := decoderAndEncoder.Encode(http.StatusOK, getTest); err != nil {
			s.logger.Error("Failed to encode data")
			jsonError.JsonError(err.Error())
			return
		}
	}
}

func (s *TestManager) CreateTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.CreateTestRequest
		json := json.New(r, s.logger, w)
		jsonError := mapjson.New(s.logger, w, r)

		err := json.DecodeAndValidationBody(&payload)
		if err != nil {
			s.logger.Error("Failed to decode body")
			jsonError.JsonError("Failed to decode body:" + err.Error())
			return
		}

		testModel := dtos.MapCreateTestRequestToModel(&payload)

		err = s.service.CreateTest(testModel)

		if err != nil {
			jsonError.JsonError("Failed to create test, error:" + err.Error())
			s.logger.Error("Failed to create test, error:" + err.Error())
			return
		}

		jsonError.JsonSuccess("Test created successfully")
	}
}
