package repository

import (
	"github.com/server/entity"
	"gorm.io/gorm"
)

type TestManagerReader interface {
	GetAllTests(user_id uint, offset, limit int) ([]entity.Test, int64, error)
	GetTestById(id uint) (*entity.Test, error)
}

type TestManagerWriter interface {
	CreateTest(data *entity.Test) error
	DeleteTest(id uint) error
	ChangeActiveStatus(status bool, testId uint) error
	IncrementCountUserPast(testId uint, count int) error
}

type TestManagerInterface interface {
	TestManagerReader
	TestManagerWriter
}

type TestManager struct {
	db *gorm.DB
}

func NewTestManager(db *gorm.DB) *TestManager {
	return &TestManager{
		db: db,
	}
}

func (s *TestManager) GetAllTests(user_id uint, lastID, limit int) ([]entity.Test, int64, error) {
	var tests []entity.Test
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

func (s *TestManager) GetTestById(id uint) (*entity.Test, error) {
	var test entity.Test

	if err := s.db.Preload("Questions.Variants").First(&test, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &test, nil
}

func (s *TestManager) CreateTest(data *entity.Test) error {
	return s.db.Create(data).Error
}

func (s *TestManager) DeleteTest(id uint) error {
	return s.db.Delete(&entity.Test{}, id).Error
}

func (s *TestManager) ChangeActiveStatus(status bool, testId uint) error {
	return s.db.Model(&entity.Test{}).
		Where("id = ?", testId).
		Update("is_active", status).Error
}

func (s *TestManager) IncrementCountUserPast(testId uint, count int) error {
	return s.db.Model(&entity.Test{}).Where("id = ?", testId).Update("count_user_past", count+1).Error
}
