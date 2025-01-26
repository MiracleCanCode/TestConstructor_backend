package dtos

import "github.com/server/internal/models"

type ValidateResultRequestPayload struct {
	Test *models.Test `json:"test" validate:"required"`
}
