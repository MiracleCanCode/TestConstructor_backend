package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/server/dtos"
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/internal/utils/json"
	"github.com/server/internal/utils/jwt"
	mapjson "github.com/server/internal/utils/mapJson"
	"github.com/server/internal/utils/middleware"
	"github.com/server/repository"
	"github.com/server/usecases"
	"go.uber.org/zap"
)

type TestManagerHandler struct {
	logger   *zap.Logger
	db       *postgresql.Db
	service  *usecases.TestManager
	userRepo *repository.User
}

func NewTestManager(logger *zap.Logger, db *postgresql.Db, router *mux.Router) {
	testManagerRepo := repository.NewTestManager(db)
	userRepo := repository.NewUser(db, logger)
	service := usecases.NewTestManager(testManagerRepo, userRepo, logger)
	handler := &TestManagerHandler{
		logger:   logger,
		db:       db,
		service:  service,
		userRepo: userRepo,
	}

	router.HandleFunc("/api/test/getById/{id}", middleware.IsAuth(handler.GetTestById())).Methods("GET")
	router.HandleFunc("/api/test/getAll", middleware.IsAuth(handler.GetAll())).Methods("POST")
	router.HandleFunc("/api/test/create", middleware.IsAuth(handler.CreateTest())).Methods("POST")
	router.HandleFunc("/api/test/delete/{id}", middleware.IsAuth(handler.DeleteTest())).Methods("DELETE")
	router.HandleFunc("/api/test/changeActive", middleware.IsAuth(handler.ChangeActiveTestStatus())).Methods("PUT")
}

func (s *TestManagerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.GetAllTestsRequest

		decoderAndEncoder := json.New(r, s.logger, w)
		jsonError := mapjson.New(s.logger, w, r)

		if err := decoderAndEncoder.DecodeAndValidationBody(&payload); err != nil {
			s.logger.Error("Failed to decode data", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			jsonError.JsonError(err.Error())
			return
		}

		getTests, count, err := s.service.GetAllTests(payload.UserId, payload.Limit, payload.Offset)
		if err != nil {
			s.logger.Error("Failed to get tests", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			jsonError.JsonError("Failed to get tests: " + err.Error())

			return
		}
		tests := dtos.SetGetAllTests(getTests, count)

		if err := decoderAndEncoder.Encode(http.StatusOK, tests); err != nil {
			s.logger.Error("Failed to encode data", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			jsonError.JsonError(err.Error())
			return
		}
	}
}

func (s *TestManagerHandler) GetTestById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		decoderAndEncoder := json.New(r, s.logger, w)
		jsonError := mapjson.New(s.logger, w, r)

		id := mux.Vars(r)["id"]
		parseId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			s.logger.Error("Failed parse id", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		userLogin, err := jwt.NewJwt(s.logger).ExtractUserFromAuthHeader(r)
		if err != nil {
			s.logger.Error("Failed extract login from user token", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}
		getTest, err := s.service.GetTestById(uint(parseId), userLogin)
		if err != nil {
			s.logger.Error("Failed to get test", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			jsonError.JsonError("Failed to get test: " + err.Error())

			return
		}

		if err := decoderAndEncoder.Encode(http.StatusOK, getTest); err != nil {
			s.logger.Error("Failed to encode data", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			jsonError.JsonError(err.Error())
			return
		}
	}
}

func (s *TestManagerHandler) CreateTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.CreateTestRequest
		json := json.New(r, s.logger, w)
		jsonError := mapjson.New(s.logger, w, r)

		if err := json.DecodeAndValidationBody(&payload); err != nil {
			s.logger.Error("Failed to decode body", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			jsonError.JsonError("Failed to decode body:" + err.Error())
			return
		}

		userLogin, err := jwt.NewJwt(s.logger).ExtractUserFromAuthHeader(r)
		if err != nil {
			s.logger.Error("Failed extract login from user token", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}
		user, err := s.userRepo.GetUserByLogin(userLogin)
		if err != nil {
			s.logger.Error("Failed to get user by login", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}
		testModel := dtos.MapCreateTestRequestToModel(&payload, user.ID)

		err = s.service.CreateTest(testModel)

		if err != nil {
			jsonError.JsonError("Failed to create test, error:")
			s.logger.Error("Failed to create test", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		jsonError.JsonSuccess("Test created successfully")
	}
}

func (s *TestManagerHandler) DeleteTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		jsonError := mapjson.New(s.logger, w, r)
		testId := mux.Vars(r)["id"]
		parseId, err := strconv.ParseUint(testId, 10, 64)
		if err != nil {
			s.logger.Error("Failed parse id", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}
		userLogin, err := jwt.NewJwt(s.logger).ExtractUserFromAuthHeader(r)
		if err != nil {
			s.logger.Error("Failed extract login from user token", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}
		test, err := s.service.GetTestById(uint(parseId), userLogin)
		if err != nil {
			s.logger.Error("Failed to get test", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			jsonError.JsonError("Failed to get test: " + err.Error())

			return
		}

		user, err := s.userRepo.GetUserByLogin(userLogin)
		if err != nil {
			s.logger.Error("Failed to get user by login", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		if test.UserID != user.ID {
			jsonError.JsonError("user can't delete this test")
			return
		}

		if err := s.service.DeleteTest(uint(parseId)); err != nil {
			s.logger.Error("Error test delete", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}
		jsonError.JsonSuccess("Success delete test")
	}
}

func (s *TestManagerHandler) ChangeActiveTestStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var payload dtos.UpdateTestActiveStatus

		json := json.New(r, s.logger, w)
		jsonError := mapjson.New(s.logger, w, r)

		userLogin, err := jwt.NewJwt(s.logger).ExtractUserFromAuthHeader(r)
		if err != nil {
			s.logger.Error("Failed to extract data from token", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		if err := json.DecodeAndValidationBody(&payload); err != nil {
			s.logger.Error("Failed to decode body", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		if err := s.service.ChangeActiveStatus(payload.IsActive, payload.TestId, userLogin); err != nil {
			s.logger.Error("Failed to change test active statuss", zap.Error(err), zap.String("method", r.Method), zap.String("endpoint", r.URL.Path))
			return
		}

		jsonError.JsonSuccess("Successed update test status")
	}
}
