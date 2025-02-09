package usecases

import (
	"errors"
	"fmt"
	"time"

	"github.com/server/internal/models"
	"github.com/server/internal/repository"
	cachemanager "github.com/server/pkg/cacheManager"
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
	testReader   repository.TestManagerReader
	testWriter   repository.TestManagerWriter
	userRepo     repository.UserReader
	logger       *zap.Logger
	cacheManager cachemanager.CacheManagerInterface
}

func NewTestManager(
	testReader repository.TestManagerReader,
	testWriter repository.TestManagerWriter,
	userRepo repository.UserReader,
	logger *zap.Logger,
) *TestManager {
	rdb := redis.New()
	cacheManager := cachemanager.New(rdb, logger)
	return &TestManager{
		testReader:   testReader,
		testWriter:   testWriter,
		userRepo:     userRepo,
		logger:       logger,
		cacheManager: cacheManager,
	}
}

func (s *TestManager) GetAllTests(userID uint, limit, offset int) ([]models.Test, int64, error) {
	cacheKey := fmt.Sprintf("tests:user:%d:limit:%d:offset:%d", userID, limit, offset)
	var result struct {
		Tests []models.Test `json:"tests"`
		Count int64         `json:"count"`
	}

	if err := s.cacheManager.Get(cacheKey, &result); err == nil {
		return result.Tests, result.Count, nil
	}

	tests, count, err := s.testReader.GetAllTests(userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	s.cacheManager.Set(cacheKey, map[string]interface{}{
		"tests": tests,
		"count": count,
	}, 10*time.Minute)

	return tests, count, nil
}

func (s *TestManager) GetTestById(id uint, userLogin string) (*models.Test, string, error) {
	cacheKey := fmt.Sprintf("test:%d", id)
	var cachedTest models.Test

	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		return nil, "", err
	}

	if err := s.cacheManager.Get(cacheKey, &cachedTest); err == nil {
		if cachedTest.UserID == user.ID {
			return &cachedTest, constants.OwnerRole, nil
		}
		return &cachedTest, constants.PassingRole, nil
	}

	test, err := s.testReader.GetTestById(id)
	if err != nil {
		return nil, "", err
	}

	if !test.IsActive && user.ID != test.UserID {
		return nil, "", errors.New("test is private")
	}

	if test.IsActive || user.ID == test.UserID {
		s.cacheManager.Set(cacheKey, test, 10*time.Minute)
	}

	if user.ID != test.UserID {
		return test, constants.PassingRole, nil
	}

	return test, constants.OwnerRole, nil
}

func (s *TestManager) CreateTest(data *models.Test) error {
	if err := s.testWriter.CreateTest(data); err != nil {
		return err
	}

	if err := s.deleteTestsFromCache(data.UserID); err != nil {
		s.logger.Warn("Failed to invalidate cache", zap.Error(err))
	}

	return nil
}

func (s *TestManager) DeleteTest(id uint, login string) error {
	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil {
		return errors.New(constants.ErrorDeleteTest)
	}

	if err := s.deleteTestsFromCache(user.ID); err != nil {
		s.logger.Warn("Failed to invalidate cache", zap.Error(err))
	}

	test, err := s.testReader.GetTestById(id)
	if err != nil {
		return errors.New(constants.GetTestByIdError)
	}

	if test.UserID != user.ID {
		return errors.New(constants.ErrorDeleteTest)
	}

	return s.testWriter.DeleteTest(id)
}

func (s *TestManager) ChangeActiveStatus(status bool, testId uint, userLogin string) error {
	test, err := s.testReader.GetTestById(testId)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		return err
	}

	if test.UserID != user.ID {
		return errors.New("user is not the author")
	}

	if err := s.deleteTestFromCache(testId); err != nil {
		s.logger.Warn("Failed to invalidate cache", zap.Error(err))
	}

	if err := s.deleteTestsFromCache(user.ID); err != nil {
		s.logger.Warn("Failed to invalidate cache", zap.Error(err))
	}

	return s.testWriter.ChangeActiveStatus(status, testId)
}

func (s *TestManager) deleteTestFromCache(testId uint) error {
	if err := s.cacheManager.Delete(fmt.Sprintf("test:%d", testId)); err != nil {
		return err
	}

	return nil
}

func (s *TestManager) deleteTestsFromCache(userId uint) error {
	if err := s.cacheManager.Delete(fmt.Sprintf("tests:user:%d:*", userId)); err != nil {
		return err
	}

	return nil
}
