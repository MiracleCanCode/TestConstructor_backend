package repository

import (
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/models"
)

type ITestManager interface {
	GetAllTests(user_id uint, offset, limit int) ([]models.Test, int64, error)
	GetTestById(id uint) (*models.Test, error)
	CreateTest(data *models.Test) error
	DeleteTest(id uint) error
	ChangeActiveStatus(status bool, testId uint) error
}

type TestManager struct {
	db *postgresql.Db
}

func NewTestManager(db *postgresql.Db) *TestManager {
	return &TestManager{
		db: db,
	}
}

func (s *TestManager) GetAllTests(user_id uint, lastID, limit int) ([]models.Test, int64, error) {
	var tests []models.Test

	query := s.db.Table("tests").Where("deleted_at IS NULL AND user_id = ?", user_id)
	if lastID > 0 {
		query = query.Where("id > ?", lastID)
	}

	if err := query.Order("id ASC").Limit(limit).Preload("Questions.Variants").Find(&tests).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := s.db.Table("tests").Where("deleted_at IS NULL AND user_id = ?", user_id).Count(&count).Error; err != nil {
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

func (s *TestManager) ChangeActiveStatus(status bool, testId uint) error {
	return s.db.Model(&models.Test{}).
		Where("id = ?", testId).
		Update("is_active", status).Error
}
