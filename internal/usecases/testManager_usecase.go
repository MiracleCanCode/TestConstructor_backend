package usecases

import (
	"errors"

	"github.com/server/internal/models"
	"github.com/server/internal/repository"
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

func (tm *TestManager) GetAllTests(userID uint, limit, offset int) ([]models.Test, int64, error) {
	return tm.testRepo.GetAllTests(userID, offset, limit)
}

func (tm *TestManager) GetTestById(id uint, userLogin string) (*models.Test, error) {
	test, err := tm.testRepo.GetTestById(id)
	if err != nil {
		tm.logger.Error("Failed to get test by id", zap.Error(err))
		return nil, err
	}

	user, err := tm.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		tm.logger.Error("Failed to get user by login", zap.Error(err))
		return nil, err
	}

	if !test.IsActive && user.ID != test.UserID {
		return nil, errors.New("test is private")
	}

	return test, nil
}

func (tm *TestManager) CreateTest(data *models.Test) error {
	return tm.testRepo.CreateTest(data)
}

func (tm *TestManager) DeleteTest(id uint) error {
	return tm.testRepo.DeleteTest(id)
}

func (tm *TestManager) ChangeActiveStatus(status bool, testId uint, userLogin string) error {
	test, err := tm.testRepo.GetTestById(testId)
	if err != nil {
		tm.logger.Error("Failed to get test", zap.Error(err))
		return err
	}

	user, err := tm.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		tm.logger.Error("Failed to get user by login", zap.Error(err))
		return err
	}

	if test.UserID != user.ID {
		return errors.New("user is not the author")
	}

	return tm.testRepo.ChangeActiveStatus(status, testId)
}
