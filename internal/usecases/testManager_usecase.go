package usecases

import (
	"fmt"
	"time"

	"github.com/server/entity"
	cachemanager "github.com/server/pkg/cacheManager"
	"github.com/server/pkg/constants"
	"github.com/server/pkg/storage/redis"
	"go.uber.org/zap"
)

type CacheManagerInterface interface {
	Get(key string, out interface{}) error
	Set(key string, value interface{}, ttl time.Duration)
	Delete(pattern string) error
}

type TestManagerRepoInterface interface {
	GetAllTests(user_id uint, offset, limit int) ([]entity.Test, int64, error)
	GetTestById(id uint) (*entity.Test, error)
	CreateTest(data *entity.Test) error
	DeleteTest(id uint) error
	ChangeActiveStatus(status bool, testId uint) error
	IncrementCountUserPast(testId uint, count int) error
}

type UserRepoInterfaceGetByLogin interface {
	GetUserByLogin(login string) (*entity.User, error)
}

type TestManager struct {
	testRepo     TestManagerRepoInterface
	userRepo     UserRepoInterfaceGetByLogin
	logger       *zap.Logger
	cacheManager CacheManagerInterface
}

func NewTestManager(
	testRepo TestManagerRepoInterface,
	userRepo UserRepoInterfaceGetByLogin,
	logger *zap.Logger,
) *TestManager {
	rdb := redis.New()
	cacheManager := cachemanager.New(rdb, logger)
	return &TestManager{
		testRepo:     testRepo,
		userRepo:     userRepo,
		logger:       logger,
		cacheManager: cacheManager,
	}
}

func (s *TestManager) GetAllTests(userID uint, limit, offset int) ([]entity.Test, int64, error) {
	cacheKey := fmt.Sprintf("tests:user:%d:limit:%d:offset:%d", userID, limit, offset)
	var result struct {
		Tests []entity.Test `json:"tests"`
		Count int64         `json:"count"`
	}

	if err := s.cacheManager.Get(cacheKey, &result); err == nil {
		return result.Tests, result.Count, nil
	}

	tests, count, err := s.testRepo.GetAllTests(userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	s.cacheManager.Set(cacheKey, map[string]interface{}{
		"tests": tests,
		"count": count,
	}, 10*time.Minute)

	return tests, count, nil
}

func (s *TestManager) GetTestById(id uint, userLogin string) (*entity.Test, string, error) {
	cacheKey := fmt.Sprintf("test:%d", id)
	var cachedTest entity.Test

	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		return nil, "", fmt.Errorf("GetTestById: failed get user by login: %w", err)
	}

	if err := s.cacheManager.Get(cacheKey, &cachedTest); err == nil {
		if cachedTest.UserID == user.ID {
			return &cachedTest, constants.OwnerRole, nil
		}
		return &cachedTest, constants.PassingRole, nil
	}

	test, err := s.testRepo.GetTestById(id)
	if err != nil {
		return nil, "", fmt.Errorf("GetTestById: failed get test by id: %w", err)
	}

	if !test.IsActive && user.ID != test.UserID {
		return nil, "", fmt.Errorf("GetTestById: %w", err)
	}

	if test.IsActive || user.ID == test.UserID {
		s.cacheManager.Set(cacheKey, test, 10*time.Minute)
	}

	if user.ID != test.UserID {
		return test, constants.PassingRole, nil
	}

	return test, constants.OwnerRole, nil
}

func (s *TestManager) CreateTest(data entity.Test) error {
	if err := s.testRepo.CreateTest(&data); err != nil {
		return fmt.Errorf("CreateTest: failed to create test: %w", err)
	}

	if err := s.deleteTestsFromCache(data.UserID); err != nil {
		return fmt.Errorf("CreateTest: failed to delete test from cache: %w", err)
	}

	return nil
}

func (s *TestManager) DeleteTest(id uint, login string) error {
	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil {
		return fmt.Errorf("DeleteTest: failed to get user by login: %w", err)
	}

	if err := s.deleteTestsFromCache(user.ID); err != nil {
		return fmt.Errorf("DeleteTest: failed to invalidate cache: %w", err)
	}

	test, err := s.testRepo.GetTestById(id)
	if err != nil {
		return fmt.Errorf("DeleteTest: failed to get test by id: %w", err)
	}

	if test.UserID != user.ID {
		return fmt.Errorf("DeleteTest: failed to delete test: %w", err)
	}

	return s.testRepo.DeleteTest(id)
}

func (s *TestManager) ChangeActiveStatus(status bool, testId uint, userLogin string) error {
	test, err := s.testRepo.GetTestById(testId)
	if err != nil {
		return fmt.Errorf("ChangeActiveStatus: failed to get test by id: %w", err)
	}

	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		return fmt.Errorf("ChangeActiveStatus: failed to get user by login:%w", err)
	}

	if test.UserID != user.ID {
		return fmt.Errorf("ChangeActiveStatus: user is not author: %w", err)
	}

	if err := s.deleteTestFromCache(testId); err != nil {
		return fmt.Errorf("ChangeActiveStatus: failed to delete test from cache: %w", err)
	}

	if err := s.deleteTestsFromCache(user.ID); err != nil {
		return fmt.Errorf("ChangeActiveStatus: failed to delete tests from cache: %w", err)
	}

	return s.testRepo.ChangeActiveStatus(status, testId)
}

func (s *TestManager) deleteTestFromCache(testId uint) error {
	if err := s.cacheManager.Delete(fmt.Sprintf("test:%d", testId)); err != nil {
		return fmt.Errorf("deleteTestFromCache: failed to delete test from cache: %w", err)
	}

	return nil
}

func (s *TestManager) deleteTestsFromCache(userId uint) error {
	if err := s.cacheManager.Delete(fmt.Sprintf("tests:user:%d:*", userId)); err != nil {
		return fmt.Errorf("deleteTestsFromCache: failed to delete tests from cache: %w", err)
	}

	return nil
}
