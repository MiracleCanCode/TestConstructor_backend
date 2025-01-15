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

func (s *TestManager) GetAllTests(user_id uint, offset, limit int) ([]models.Test, int64, error) {
	var (
		tests []models.Test
		count int64
	)

	if err := s.db.Table("tests").Where("deleted_at is null and user_id = ?", user_id).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Table("tests").Where("user_id = ? AND deleted_at is null", &user_id).
		Order("id ASC").Offset(offset).Limit(limit).
		Preload("Questions.Variants").Find(&tests).Error; err != nil {
		return nil, 0, err
	}

	return tests, count, nil
}

func (s *TestManager) GetTestById(id uint) (*models.Test, error) {
	var test models.Test

	if err := s.db.Preload("Questions.Variants").First(&test, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &test, nil
}

func (s *TestManager) CreateTest(data *models.Test) error {
	return s.db.Create(data).Error
}

func (s *TestManager) DeleteTest(id uint) error {
	return s.db.Delete(&models.Test{}, id).Error
}
