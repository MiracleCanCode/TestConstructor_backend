package validateresulttest

import "github.com/server/models"

type ValidateResultTestRequest struct {
	Test *models.Test `json:"test" validate:"required"`
}
