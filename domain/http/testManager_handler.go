package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/internal/usecases"
	"github.com/server/pkg/constants"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/errors"
	"github.com/server/pkg/json"
	"github.com/server/pkg/jwt"
	"github.com/server/pkg/middleware"
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

	router.HandleFunc("/api/test/getById/{id}", middleware.IsAuth(handler.GetTestById())).Methods(http.MethodGet)
	router.HandleFunc("/api/test/getAll", middleware.IsAuth(handler.GetAll())).Methods(http.MethodPost)
	router.HandleFunc("/api/test/create", middleware.IsAuth(handler.CreateTest())).Methods(http.MethodPost)
	router.HandleFunc("/api/test/delete/{id}", middleware.IsAuth(handler.DeleteTest())).Methods(http.MethodDelete)
	router.HandleFunc("/api/test/changeActive", middleware.IsAuth(handler.ChangeActiveTestStatus())).Methods(http.MethodPut)
}

func (s *TestManagerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.GetAllTestsRequest
		decoderAndEncoder := json.New(r, s.logger, w)

		if err := decoderAndEncoder.DecodeAndValidationBody(&payload); err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to decode data", constants.InternalServerError)
			return
		}

		getTests, count, err := s.service.GetAllTests(payload.UserId, payload.Limit, payload.Offset)
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to get tests", constants.ErrorGetAllTests)
			return
		}
		tests := dtos.SetGetAllTests(getTests, count)

		if err := decoderAndEncoder.Encode(http.StatusOK, tests); err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to encode data", constants.InternalServerError)
			return
		}
	}
}

func (s *TestManagerHandler) GetTestById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		decoderAndEncoder := json.New(r, s.logger, w)
		id := mux.Vars(r)["id"]
		parseId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to parse id", constants.InternalServerError)
			return
		}

		userLogin, err := jwt.NewJwt(s.logger).ExtractUserFromCookie(r, "token")
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to extract login from user token", constants.InternalServerError)
			return
		}

		getTest, role, err := s.service.GetTestById(uint(parseId), userLogin)
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to get test", constants.GetTestByIdError)
			return
		}

		res := dtos.MapTestModelToGetTestByIdResponse(getTest, role)
		if err := decoderAndEncoder.Encode(http.StatusOK, res); err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to encode data", constants.InternalServerError)
			return
		}
	}
}

func (s *TestManagerHandler) CreateTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload dtos.CreateTestRequest
		json := json.New(r, s.logger, w)

		if err := json.DecodeAndValidationBody(&payload); err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to decode body", constants.InternalServerError)
			return
		}

		userLogin, err := jwt.NewJwt(s.logger).ExtractUserFromCookie(r, "token")
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to extract login from user token", constants.InternalServerError)
			return
		}

		user, err := s.userRepo.GetUserByLogin(userLogin)
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to get user by login", constants.NotFoundUser)
			return
		}
		testModel := dtos.MapCreateTestRequestToModel(&payload, user.ID)

		err = s.service.CreateTest(testModel)
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to create test", constants.ErrorCreateTest)
			return
		}

	}
}

func (s *TestManagerHandler) DeleteTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		testId := mux.Vars(r)["id"]

		parseId, err := strconv.ParseUint(testId, 10, 64)
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed parse id", constants.InternalServerError)
			return
		}

		login, err := jwt.NewJwt(s.logger).ExtractUserFromCookie(r, "token")
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed extract login from user token", constants.InternalServerError)
			return
		}

		if err := s.service.DeleteTest(uint(parseId), login); err != nil {
			errors.HandleError(s.logger, w, r, err, "Error test delete", constants.ErrorDeleteTest)
			return
		}

	}
}

func (s *TestManagerHandler) ChangeActiveTestStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var payload dtos.UpdateTestActiveStatus

		json := json.New(r, s.logger, w)

		userLogin, err := jwt.NewJwt(s.logger).ExtractUserFromCookie(r, "token")
		if err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to extract data from token", constants.InternalServerError)
			return
		}

		if err := json.DecodeAndValidationBody(&payload); err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to decode body", constants.InternalServerError)
			return
		}

		if err := s.service.ChangeActiveStatus(payload.IsActive, payload.TestId, userLogin); err != nil {
			errors.HandleError(s.logger, w, r, err, "Failed to change test active status", constants.ErrorChangeActiveTest)
			return
		}

	}
}
