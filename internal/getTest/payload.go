package getTest

import "github.com/server/models"

func SetDataToGetAllTestsResponse(tests []*models.Test, count int64) *GetAllTestsResponse {
	return &GetAllTestsResponse{
		Tests: tests,
		Count: count,
	}
}

type GetAllTestsRequest struct {
	Login  string `json:"login" validate:"required"`
	Limit  int    `json:"limit" validate:"gte=0"`
	Offset int    `json:"offset" validate:"gte=0"`
}

type GetTestByIdRequest struct {
	TestId uint `json:"test_id" validate:"required"`
}

type GetAllTestsResponse struct {
	Tests []*models.Test `json:"tests"`
	Count int64          `json:"count"`
}

type GetTestByIdResponse struct {
	models.Test
}
