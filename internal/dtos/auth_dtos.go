package dtos

import "github.com/server/entity"

func (r *RegistrationRequest) ToUser() entity.User {
	return entity.User{
		Name:     r.Name,
		Login:    r.Login,
		Password: r.Password,
		Avatar:   r.Avatar,
		Email:    r.Email,
	}
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *entity.User `json:"user"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegistrationResponse struct {
	Name         string  `json:"name"`
	Login        string  `json:"login"`
	Avatar       *string `json:"avatar"`
	Email        string  `json:"email"`
	RefreshToken string  `json:"refresh_token"`
}

type RegistrationRequest struct {
	Name     string  `json:"name" validate:"required"`
	Avatar   *string `json:"avatar,omitempty"`
	Login    string  `json:"login" validate:"required"`
	Password string  `json:"password" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
}
