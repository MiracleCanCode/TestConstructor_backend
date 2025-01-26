package dtos

import "github.com/server/internal/models"

func (r *RegistrationRequest) ToUser() models.User {
	return models.User{
		Name:     r.Name,
		Login:    r.Login,
		Password: r.Password,
		Avatar:   r.Avatar,
		Email:    r.Email,
	}
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegistrationResponse struct {
	Name         string  `json:"name"`
	Login        string  `json:"login"  `
	Password     string  `json:"password"`
	Avatar       *string `json:"avatar"`
	Email        string  `json:"email" `
	RefreshToken string  `json:"refresh_token" `
}

type RegistrationRequest struct {
	Name     string  `json:"name" validate:"required"`
	Avatar   *string `json:"avatar,omitempty"`
	Login    string  `json:"login" validate:"required"`
	Password string  `json:"password" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
}
