package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/server/internal/models"
	"github.com/server/internal/repository"
	"github.com/server/pkg/constants"
	"github.com/server/pkg/db/redis"
	"go.uber.org/zap"
)

type TestManagerInterface interface {
	GetAllTests(userID uint, limit, offset int) ([]models.Test, int64, error)
	GetTestById(id uint, userLogin string) (*models.Test, error)
	CreateTest(data *models.Test) error
	DeleteTest(id uint) error
	ChangeActiveStatus(status bool, testId uint, userLogin string) error
}

type TestManager struct {
	testRepo repository.TestManagerInterface
	userRepo repository.UserInterface
	logger   *zap.Logger
}

func NewTestManager(
	testRepo repository.TestManagerInterface,
	userRepo repository.UserInterface,
	logger *zap.Logger,
) *TestManager {
	return &TestManager{
		testRepo: testRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *TestManager) GetAllTests(userID uint, limit, offset int) ([]models.Test, int64, error) {
	rdb := redis.New()
	cacheKey := fmt.Sprintf("tests:user:%d:limit:%d:offset:%d", userID, limit, offset)

	cachedData, err := rdb.Get(cacheKey)
	if err == nil {
		var result struct {
			Tests []models.Test `json:"tests"`
			Count int64         `json:"count"`
		}

		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			return result.Tests, result.Count, nil
		}
		s.logger.Error("Failed to unmarshal cache", zap.Error(err))
	}

	tests, count, err := s.testRepo.GetAllTests(userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	cacheValue, err := json.Marshal(map[string]interface{}{
		"tests": tests,
		"count": count,
	})
	if err == nil {
		_ = rdb.Set(cacheKey, cacheValue, 10*time.Minute)
	} else {
		s.logger.Warn("Failed to marshal cache data", zap.Error(err))
	}

	return tests, count, nil
}

func (s *TestManager) GetTestById(id uint, userLogin string) (*models.Test, string, error) {
	test, err := s.testRepo.GetTestById(id)
	if err != nil {
		s.logger.Error("Failed to get test by id", zap.Error(err))
		return nil, "", err
	}

	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		s.logger.Error("Failed to get user by login", zap.Error(err))
		return nil, "", err
	}

	if !test.IsActive && user.ID != test.UserID {
		return nil, "", errors.New("test is private")
	}

	if user.ID != test.UserID {
		return test, constants.PassingRole, nil
	}

	return test, constants.OwnerRole, nil
}

func (s *TestManager) CreateTest(data *models.Test) error {
	return s.testRepo.CreateTest(data)
}

func (s *TestManager) DeleteTest(id uint, login string) error {
	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil {
		return errors.New(constants.ErrorDeleteTest)
	}

	test, _, err := s.GetTestById(id, login)
	if err != nil {
		return errors.New(constants.GetTestByIdError)
	}

	if test.UserID != user.ID {
		return errors.New(constants.ErrorDeleteTest)
	}

	return s.testRepo.DeleteTest(id)
}

func (s *TestManager) ChangeActiveStatus(status bool, testId uint, userLogin string) error {
	test, err := s.testRepo.GetTestById(testId)
	if err != nil {
		s.logger.Error("Failed to get test", zap.Error(err))
		return err
	}

	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		s.logger.Error("Failed to get user by login", zap.Error(err))
		return err
	}

	if test.UserID != user.ID {
		return errors.New("user is not the author")
	}

	return s.testRepo.ChangeActiveStatus(status, testId)
}
