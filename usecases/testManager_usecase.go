package usecases

import (
	"github.com/server/models"
	"github.com/server/repository"

	"github.com/server/internal/utils/db/postgresql"
	"go.uber.org/zap"
)

type TestManager struct {
	db         *postgresql.Db
	logger     *zap.Logger
	repository *repository.TestManager
}

func NewTestManager(db *postgresql.Db, logger *zap.Logger) *TestManager {
	return &TestManager{
		db:         db,
		logger:     logger,
		repository: repository.NewTestManager(db),
	}
}

func (s *TestManager) GetAllTests(user_id uint, limit, offset int) ([]models.Test, int64, error) {
	return s.repository.GetAllTests(user_id, offset, limit)
}

func (s *TestManager) GetTestById(id uint) (*models.Test, error) {
	return s.repository.GetTestById(id)
}
func (s *TestManager) CreateTest(data *models.Test) error {
	return s.repository.CreateTest(data)
}

func (s *TestManager) DeleteTest(id uint) error {
	return s.repository.DeleteTest(id)
}
