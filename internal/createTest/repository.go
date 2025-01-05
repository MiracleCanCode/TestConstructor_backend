package createTest

import (
	"github.com/server/models"
	"github.com/server/pkg/db/postgresql"
)

type IRepository interface {
	Create(data *models.Test) error
}

type Repository struct {
	db *postgresql.Db
}

func NewRepository(db *postgresql.Db) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) Create(data *models.Test) error {
	res := s.db.Create(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
