package repository

import (
	"github.com/server/dtos"
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IUser interface {
	Update(user *dtos.UpdateUserRequest) error
	Create(user models.User) error
	GetByLogin(login string) (*models.User, error)
}

type User struct {
	db     *postgresql.Db
	logger *zap.Logger
}

func NewUser(db *postgresql.Db, logger *zap.Logger) *User {
	return &User{
		db:     db,
		logger: logger,
	}
}

func (s *User) UpdateUser(user *dtos.UpdateUserRequest) error {
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

func (s *User) CreateUser(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	if err := s.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *User) GetUserByLogin(login string) (*models.User, error) {
	var user models.User

	if err := s.db.Where("login = ?", login).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *User) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}