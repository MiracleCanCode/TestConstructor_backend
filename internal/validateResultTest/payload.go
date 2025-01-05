package validateresulttest

import "github.com/server/models"

type RequestPayload struct {
	Test *models.Test `json:"test" validate:"required"`
}
