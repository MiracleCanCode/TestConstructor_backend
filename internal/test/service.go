package test

import (
	"github.com/server/models"

	"github.com/server/pkg/db/postgresql"
	"go.uber.org/zap"
)

type Service struct {
	db         *postgresql.Db
	logger     *zap.Logger
	repository *Repository
}

func NewService(db *postgresql.Db, logger *zap.Logger, repository *Repository) *Service {
	if repository == nil {
		repository = NewRepository(db)
	}
	return &Service{
		db:         db,
		logger:     logger,
		repository: repository,
	}
}

func (s *Service) GetAll(login string, limit, offset int) ([]models.Test, int64, error) {
	return s.repository.GetAll(login, offset, limit)
}

func (s *Service) GetById(id uint) (*models.Test, error) {
	return s.repository.GetById(id)
}
func (s *Service) Create(data *models.Test) error {
	createTest := s.repository.Create(data)
	if createTest != nil {
		return createTest
	}

	return nil
}
