package user

import (
	"github.com/server/models"
	"github.com/server/pkg/db/postgresql"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IRepository interface {
	Update(user *UpdateRequest) error
	Create(user models.User) error
	GetByLogin(login string) (*models.User, error)
}

type Repository struct {
	db *postgresql.Db
	logger *zap.Logger
}

func NewRepository(db *postgresql.Db, logger *zap.Logger) *Repository {
	return &Repository{
		db: db,
		logger: logger,
	}
}

func (s *Repository) Update(user *UpdateRequest) error {
	updateData := map[string]interface{}{}

	if user.Data.Name != nil {
		updateData["name"] = &user.Data.Name
	}
	if user.Data.Avatar != nil {
		updateData["avatar"] = &user.Data.Avatar
	}

	if len(updateData) == 0 {
		return nil 
	}

	if err := s.db.Model(&models.User{}).Where("login = ?", user.UserLogin).Updates(updateData).Error; err != nil {
		s.logger.Error("Failed update user data, error:" + err.Error())
		return err
	}

	return nil
}

func (s *Repository) Create(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	result := s.db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Repository) GetByLogin(login string) (*models.User, error) {

	var user models.User
	result := s.db.Where("login = ?", login).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
