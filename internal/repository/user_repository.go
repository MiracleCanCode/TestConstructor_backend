package repository

import (
	"fmt"

	"github.com/server/entity"
	"github.com/server/internal/dtos"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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
		return fmt.Errorf("UpdateUser: failed to update user data: %w", err)
	}

	return nil
}

func (s *User) CreateUser(user entity.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("CreateUser: failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)
	s.logger.Info(user.RefreshToken)
	if err := s.db.Create(&user).Error; err != nil {
		return fmt.Errorf("CreateUser: failed to create user: %w", err)
	}

	return nil
}

func (s *User) GetUserByLogin(login string) (*entity.User, error) {
	var user entity.User

	if err := s.db.Select("id, login, email, name, avatar, password, refresh_token").
		Where("login = ?", login).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("GetUserByLogin: failed to get user by login: %w", err)
	}

	return &user, nil
}

func (s *User) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User

	if err := s.db.Select("id, login, email, name, avatar, refresh_token").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("GetUserByEmail: failed to get user by email: %w", err)
	}

	return &user, nil
}

func (s *User) DeleteRefreshToken(login string) error {
	var user entity.User

	if err := s.db.Where("login = ?", login).First(&user).Error; err != nil {
		return fmt.Errorf("DeleteRefreshToken: failed to get user by login: %w", err)
	}

	user.RefreshToken = ""

	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("DeleteRefreshToken: failed save user: %w", err)
	}

	return nil
}
