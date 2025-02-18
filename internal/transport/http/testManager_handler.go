package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/server/entity"
	"github.com/server/internal/dtos"
	"github.com/server/internal/repository"
	"github.com/server/internal/transport/http/middleware"
	"github.com/server/internal/usecases"
	"github.com/server/pkg/constants"
	errorshandler "github.com/server/pkg/errorsHandler"
	"github.com/server/pkg/json"
	"github.com/server/pkg/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	GetUserByLogin(login string) (*entity.User, error)
}

type TestMangerUseCaseInterface interface {
	GetAllTests(userID uint, limit, offset int) ([]entity.Test, int64, error)
	GetTestById(id uint, userLogin string) (*entity.Test, string, error)
	CreateTest(data entity.Test) error
	DeleteTest(id uint, login string) error
	ChangeActiveStatus(status bool, testId uint, userLogin string) error
}

type TestManagerHandler struct {
	logger   *zap.Logger
	db       *gorm.DB
	service  TestMangerUseCaseInterface
	userRepo UserRepoInterface
}

func NewTestManagerHandler(logger *zap.Logger, db *gorm.DB, router *mux.Router) {
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
		errors := errorshandler.New(s.logger, w, r)
		var payload dtos.GetAllTestsRequest
		decoderAndEncoder := json.New(r, s.logger, w)

		if err := decoderAndEncoder.Decode(&payload); err != nil {
			errors.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		getTests, count, err := s.service.GetAllTests(payload.UserId, payload.Limit, payload.Offset)
		if err != nil {
			errors.HandleError(constants.ErrorGetAllTests, http.StatusNotFound, err)
			return
		}
		tests := dtos.SetGetAllTests(getTests, count)

		if err := decoderAndEncoder.Encode(http.StatusOK, tests); err != nil {
			errors.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *TestManagerHandler) GetTestById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		errors := errorshandler.New(s.logger, w, r)
		decoderAndEncoder := json.New(r, s.logger, w)
		id := mux.Vars(r)["id"]
		parseId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			s.logger.Error("GetTestById: failed parse test id", zap.Error(err))
			errors.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		userLogin, err := jwt.NewJwt(s.logger).ExtractUserFromToken(r)
		if err != nil {
			s.logger.Error("GetTestById: failed extract user from token", zap.Error(err))
			errors.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		getTest, role, err := s.service.GetTestById(uint(parseId), userLogin)
		if err != nil {
			s.logger.Error("GetTestById: failed get test by id", zap.Error(err))
			errors.HandleError(constants.GetTestByIdError, http.StatusNotFound, err)
			return
		}

		res := dtos.MapTestToGetTestResponse(getTest, role, getTest.UserID)
		if err := decoderAndEncoder.Encode(http.StatusOK, res); err != nil {
			s.logger.Error("GetTestById: failed encode response body", zap.Error(err))
			errors.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *TestManagerHandler) CreateTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		errors := errorshandler.New(s.logger, w, r)
		var payload dtos.CreateTestRequest
		json := json.New(r, s.logger, w)

		if err := json.Decode(&payload); err != nil {
			s.logger.Error("CreateTest: failed decode request body", zap.Error(err))
			errors.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		userLogin, err := jwt.NewJwt(s.logger).ExtractUserFromToken(r)
		if err != nil {
			s.logger.Error("CreateTest: failed extract user from token", zap.Error(err))
			errors.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		user, err := s.userRepo.GetUserByLogin(userLogin)
		if err != nil {
			s.logger.Error("CreateTest: failed get user by login", zap.Error(err))
			errors.HandleError(constants.NotFoundUser, http.StatusNotFound, err)
			return
		}
		testModel := dtos.MapCreateTestRequestToModel(&payload, user.ID)

		err = s.service.CreateTest(testModel)
		if err != nil {
			s.logger.Error("CreateTest: failed create test", zap.Error(err))
			errors.HandleError(constants.ErrorCreateTest, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *TestManagerHandler) DeleteTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		errors := errorshandler.New(s.logger, w, r)
		testId := mux.Vars(r)["id"]
		parseId, err := strconv.ParseUint(testId, 10, 64)
		if err != nil {
			s.logger.Error("DeleteTest: failed parse test id", zap.Error(err))
			errors.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		login, err := jwt.NewJwt(s.logger).ExtractUserFromToken(r)
		if err != nil {
			s.logger.Error("DeleteTest: failed extract user from token", zap.Error(err))
			errors.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		if err := s.service.DeleteTest(uint(parseId), login); err != nil {
			s.logger.Error("DeleteTest: failed delete test", zap.Error(err))
			errors.HandleError(constants.ErrorDeleteTest, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *TestManagerHandler) ChangeActiveTestStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		errors := errorshandler.New(s.logger, w, r)
		var payload dtos.UpdateTestActiveStatus
		json := json.New(r, s.logger, w)

		userLogin, err := jwt.NewJwt(s.logger).ExtractUserFromToken(r)
		if err != nil {
			s.logger.Error("ChangeActiveTestStatus: failed extract user from token", zap.Error(err))
			errors.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		if err := json.Decode(&payload); err != nil {
			s.logger.Error("ChangeActiveTestStatus: failed decode request body", zap.Error(err))
			errors.HandleError(constants.InternalServerError, http.StatusInternalServerError, err)
			return
		}

		if err := s.service.ChangeActiveStatus(payload.IsActive, payload.TestId, userLogin); err != nil {
			s.logger.Error("ChangeActiveTestStatus: failed change active test status", zap.Error(err))
			errors.HandleError(constants.ErrorChangeActiveTest, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
