package dtos

type UpdateUserRequest struct {
	UserLogin string `json:"user_login" validate:"required"`
	Data      struct {
		Name   *string `json:"name"`
		Avatar *string `json:"avatar"`
	} `json:"data" validate:"required"`
}