package getTest

import "github.com/server/models"

type GetAllTestsRequest struct {
	UserId string `json:"user_id" validate:"required"`
}

type GetTestByIdRequest struct {
	TestId string `json:"test_id" validate:"required"`
}

type GetAllTestsResponse struct {
	tests []models.Test `json:"tests"`
}

type GetTestByIdResponse struct {
	models.Test
}
