package repository

import (
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/models"
)

type ITestManager interface {
	GetAll(login string, offset, limit int) ([]models.Test, int64, error)
	GetById(id uint) (*models.Test, error)
}

type TestManager struct {
	db *postgresql.Db
}

func NewTestManager(db *postgresql.Db) *TestManager {
	return &TestManager{
		db: db,
	}
}

func (s *TestManager) GetAllTests(login string, offset, limit int) ([]models.Test, int64, error) {
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

func (s *TestManager) GetTestById(id uint) (*models.Test, error) {
	var test models.Test

	err := s.db.Preload("Questions.Variants").First(&test, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &test, nil
}

func (s *TestManager) CreateTest(data *models.Test) error {
	res := s.db.Create(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
