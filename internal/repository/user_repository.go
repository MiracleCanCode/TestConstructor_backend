package repository

import (
	"github.com/server/entity"
	"github.com/server/internal/dtos"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserReader interface {
	GetUserByLogin(login string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
}

type UserWriter interface {
	UpdateUser(user *dtos.UpdateUserRequest) error
	CreateUser(user entity.User) error
	DeleteRefreshToken(login string) error
}

type UserInterface interface {
	UserReader
	UserWriter
}

type User struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUser(db *gorm.DB, logger *zap.Logger) *User {
	return &User{
		db:     db,
		logger: logger,
	}
}

func (s *User) UpdateUser(user *dtos.UpdateUserRequest) error {
	updateData := make(map[string]interface{})

	if user.Data.Name != nil {
		updateData["name"] = *user.Data.Name
	}
	if user.Data.Avatar != nil {
		updateData["avatar"] = *user.Data.Avatar
	}

	if len(updateData) == 0 {
		return nil
	}

	if err := s.db.Model(&entity.User{}).
		Where("login = ?", user.UserLogin).
		Updates(updateData).Error; err != nil {
		s.logger.Error("Failed to update user data", zap.Error(err))
		return err
	}

	return nil
}

func (s *User) CreateUser(user entity.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	s.logger.Info(user.RefreshToken)
	if err := s.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *User) GetUserByLogin(login string) (*entity.User, error) {
	var user entity.User

	if err := s.db.Select("id, login, email, name, avatar, password, refresh_token").
		Where("login = ?", login).
		First(&user).Error; err != nil {
		s.logger.Error("Failed to get user by login", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (s *User) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User

	if err := s.db.Select("id, login, email, name, avatar, refresh_token").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		s.logger.Error("Failed to get user by email", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (s *User) DeleteRefreshToken(login string) error {
	var user entity.User

	if err := s.db.Where("login = ?", login).First(&user).Error; err != nil {
		return err
	}

	user.RefreshToken = ""

	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
