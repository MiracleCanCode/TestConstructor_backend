package usecases

import (
	"errors"

	"github.com/server/models"
	"github.com/server/repository"

	"github.com/server/internal/utils/db/postgresql"
	"go.uber.org/zap"
)

type TestManager struct {
	db         *postgresql.Db
	logger     *zap.Logger
	repository *repository.TestManager
	userRepo   *repository.User
}

func NewTestManager(db *postgresql.Db, logger *zap.Logger) *TestManager {
	return &TestManager{
		db:         db,
		logger:     logger,
		repository: repository.NewTestManager(db),
		userRepo:   repository.NewUser(db, logger),
	}
}

func (s *TestManager) GetAllTests(user_id uint, limit, offset int) ([]models.Test, int64, error) {
	return s.repository.GetAllTests(user_id, offset, limit)
}

func (s *TestManager) GetTestById(id uint, userLogin string) (*models.Test, error) {
	test, err := s.repository.GetTestById(id)
	if err != nil {
		s.logger.Error("Failed to get test by id", zap.Error(err))
		return nil, err
	}

	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		s.logger.Error("Failed to get user by login", zap.Error(err))
		return nil, err
	}

	if !test.IsActive {
		if user.ID != test.UserID {
			return nil, errors.New("test is private")
		}
	}

	return test, nil

}
func (s *TestManager) CreateTest(data *models.Test) error {
	return s.repository.CreateTest(data)
}

func (s *TestManager) DeleteTest(id uint) error {
	return s.repository.DeleteTest(id)
}

func (s *TestManager) ChangeActiveStatus(status bool, testId uint, userLogin string) error {
	test, err := s.repository.GetTestById(testId)
	if err != nil {
		s.logger.Error("Failed get test", zap.Error(err))
		return err
	}

	user, err := s.userRepo.GetUserByLogin(userLogin)
	if err != nil {
		s.logger.Error("Failed to get user by login", zap.Error(err))
		return err
	}

	if test.UserID != user.ID {
		return errors.New("user is not author")
	}

	return s.repository.ChangeActiveStatus(status, testId)
}
