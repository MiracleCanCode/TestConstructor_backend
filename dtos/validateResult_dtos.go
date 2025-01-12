package dtos

import "github.com/server/models"

type ValidateResultRequestPayload struct {
	Test *models.Test `json:"test" validate:"required"`
}
