package repository

import (
	"fmt"

	"github.com/server/entity"
	"gorm.io/gorm"
)

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
		return nil, 0, fmt.Errorf("GetAllTests: failed to get all tests: %w", err)
	}

	var count int64
	if err := s.db.Table("tests").Where("deleted_at IS NULL AND user_id = ?", user_id).Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("GetAllTests: failed to get count: %w", err)
	}

	return tests, count, nil
}

func (s *TestManager) GetTestById(id uint) (*entity.Test, error) {
	var test entity.Test

	if err := s.db.Preload("Questions.Variants").First(&test, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("GetTestById: failed to get test by id: %w", err)
	}

	return &test, nil
}

func (s *TestManager) CreateTest(data *entity.Test) error {
	if err := s.db.Create(data).Error; err != nil {
		return fmt.Errorf("CreateTest: failed to create test: %w", err)
	}
	return nil
}

func (s *TestManager) DeleteTest(id uint) error {
	if err := s.db.Delete(&entity.Test{}, id).Error; err != nil {
		return fmt.Errorf("DeleteTest: failed delete test: %w", err)
	}
	return nil
}

func (s *TestManager) ChangeActiveStatus(status bool, testId uint) error {
	if err := s.db.Model(&entity.Test{}).
		Where("id = ?", testId).
		Update("is_active", status).Error; err != nil {
		return fmt.Errorf("ChangeActiveStatus: failed to change active status test: %w", err)
	}
	return nil
}

func (s *TestManager) IncrementCountUserPast(testId uint, count int) error {
	if err := s.db.Model(&entity.Test{}).Where("id = ?", testId).Update("count_user_past", count+1).Error; err != nil {
		return fmt.Errorf("IncrementCountUserPast: failed to increment count user past: %w", err)
	}

	return nil
}
