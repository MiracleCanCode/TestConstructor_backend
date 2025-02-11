package dtos

import "github.com/server/entity"

type ValidateResultRequestPayload struct {
	Test *entity.Test `json:"test" validate:"required"`
}
