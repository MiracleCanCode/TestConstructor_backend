package getTest

import (
	"github.com/server/models"
	"github.com/server/pkg/db"
	"go.uber.org/zap"
)

type GetTestService struct {
	db         *db.Db
	logger     *zap.Logger
	repository *GetTestRepository
}

func NewGetTestService(db *db.Db, logger *zap.Logger, repository *GetTestRepository) *GetTestService {
	if repository == nil {
		repository = NewGetTestRepository(db)
	}
	return &GetTestService{
		db:         db,
		logger:     logger,
		repository: repository,
	}
}

func (s *GetTestService) GetAllTests(login string, limit, offset int) ([]models.Test, int64, error) {
	return s.repository.GetAllTests(login, offset, limit)
}

func (s *GetTestService) GetTestById(id uint) (*models.Test, error) {
	return s.repository.GetTestById(id)
}
