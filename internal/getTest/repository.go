package getTest

import (
	"github.com/server/models"
	"github.com/server/pkg/db"
)

type GetTestRepository struct {
	db *db.Db
}

func NewGetTestRepository(db *db.Db) *GetTestRepository {
	return &GetTestRepository{
		db: db,
	}
}

func (s *GetTestRepository) GetAllTests(login string, offset, limit int) ([]models.Test, int64, error) {
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

func (s *GetTestRepository) GetTestById(id string) (*models.Test, error) {
	var test models.Test

	err := s.db.Preload("Questions.Variants").First(&test, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &test, nil
}
