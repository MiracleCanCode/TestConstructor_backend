package dtos

import "github.com/server/internal/models"

func ToGetUserByLoginResponse(user *models.User) *GetUserByLoginResponse {
	return &GetUserByLoginResponse{
		Login:  user.Login,
		Name:   user.Name,
		Avatar: user.Avatar,
		Id:     user.ID,
		Email:  user.Email,
	}
}

type UpdateUserRequest struct {
	UserLogin string `json:"user_login" validate:"required"`
	Data      struct {
		Name   *string `json:"name"`
		Avatar *string `json:"avatar"`
	} `json:"data" validate:"required"`
}

type GetUserByLoginRequest struct {
	Login string `json:"login" validate:"required"`
}

type GetUserByLoginResponse struct {
	Login  string  `json:"login"`
	Email  string  `json:"email"`
	Id     uint    `json:"id"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
}
