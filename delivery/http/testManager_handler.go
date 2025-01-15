package http

import (
	"net/http"
	"strconv"
	"strings"

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

type TestManager struct {
	logger   *zap.Logger
	db       *postgresql.Db
	service  *usecases.TestManager
	userRepo *repository.User
}

func NewTestManager(logger *zap.Logger, db *postgresql.Db, router *mux.Router) {
	handler := &TestManager{
		logger:   logger,
		db:       db,
		service:  usecases.NewTestManager(db, logger),
		userRepo: repository.NewUser(db, logger),
	}

	router.HandleFunc("/api/test/getById/{id}", handler.GetTestById()).Methods("GET")
	router.HandleFunc("/api/test/getAll", middleware.IsAuth(handler.GetAll())).Methods("POST")
	router.HandleFunc("/api/test/create", handler.CreateTest()).Methods("POST")
	router.HandleFunc("/api/test/delete/{id}", middleware.IsAuth(handler.DeleteTest())).Methods("DELETE")
}

func (s *TestManager) GetAll() http.HandlerFunc {
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

func (s *TestManager) GetTestById() http.HandlerFunc {
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
		getTest, err := s.service.GetTestById(uint(parseId))
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

func (s *TestManager) CreateTest() http.HandlerFunc {
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

		authHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, claims, err := jwt.NewJwt("SUPERSECRETKEYFORBESTAPPINTHEWORLD").VerifyToken(tokenString)
		if err != nil {
			return
		}

		userLogin := claims["login"].(string)
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

func (s *TestManager) DeleteTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
	}
}
