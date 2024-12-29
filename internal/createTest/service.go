package createTest

import (
	"github.com/server/models"
	"github.com/server/pkg/db"
	"go.uber.org/zap"
)

type CreateTestService struct {
	db         *db.Db
	logger     *zap.Logger
	repository *CreateTestRepository
}

func NewCreateTestService(db *db.Db, logger *zap.Logger) *CreateTestService {
	return &CreateTestService{
		db:         db,
		logger:     logger,
		repository: NewCreateTestRepository(db),
	}
}

func (s *CreateTestService) CreateTest(data *models.Test) error {
	createTest := s.repository.Create(data)
	if createTest != nil {
		return createTest
	}

	return nil
}
