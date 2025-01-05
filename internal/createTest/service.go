package createTest

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

func NewService(db *postgresql.Db, logger *zap.Logger) *Service {
	return &Service{
		db:         db,
		logger:     logger,
		repository: NewRepository(db),
	}
}

func (s *Service) Create(data *models.Test) error {
	createTest := s.repository.Create(data)
	if createTest != nil {
		return createTest
	}

	return nil
}
