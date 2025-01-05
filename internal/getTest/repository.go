package getTest

import (
	"github.com/server/models"
	"github.com/server/pkg/db/postgresql"
)

type IRepository interface {
	GetAll(login string, offset, limit int) ([]models.Test, int64, error)
	GetById(id uint) (*models.Test, error)
}

type Repository struct {
	db *postgresql.Db
}

func NewRepository(db *postgresql.Db) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) GetAll(login string, offset, limit int) ([]models.Test, int64, error) {
	var (
		tests []models.Test
		count int64
	)

	err := s.db.Table("tests").Where("deleted_at is null and author_login = ?", login).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Table("tests").Where("author_login = ? AND deleted_at is null", &login).
		Order("id ASC").Offset(offset).Limit(limit).
		Preload("Questions.Variants").Find(&tests).Error
	if err != nil {
		return nil, 0, err
	}
	return tests, count, nil
}

func (s *Repository) GetById(id uint) (*models.Test, error) {
	var test models.Test

	err := s.db.Preload("Questions.Variants").First(&test, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &test, nil
}
