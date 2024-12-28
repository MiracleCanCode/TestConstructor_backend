package createTest

import (
	"github.com/server/models"
	"github.com/server/pkg/db"
)

type CreateTestRepository struct {
	db *db.Db
}

func NewCreateTestRepository(db *db.Db) *CreateTestRepository {
	return &CreateTestRepository{
		db: db,
	}
}

func (s *CreateTestRepository) Create(data *models.Test) error {
	res := s.db.Create(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
