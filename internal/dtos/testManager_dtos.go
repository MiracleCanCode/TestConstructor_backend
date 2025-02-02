package dtos

import "github.com/server/internal/models"

func MapTestModelToGetTestByIdResponse(req *models.Test, userRole string) *GetTestByIdResponse {
	return &GetTestByIdResponse{
		Test: *req,
		Role: userRole,
	}
}

func MapCreateTestRequestToModel(req *CreateTestRequest, userId uint) *models.Test {
	test := &models.Test{
		Name:      req.Name,
		UserID:    userId,
		Questions: mapQuestions(req.Questions),
	}

	return test
}

func MapCreateAnonymusTestRequestToModel(req *CreateAnonymusTestRequest) *models.Test {
	test := &models.Test{
		Name:      req.Name,
		Questions: mapQuestions(req.Questions),
	}

	return test
}

func mapQuestions(questions []CreateQuestionInput) []models.Question {
	mappedQuestions := make([]models.Question, len(questions))

	for i, question := range questions {
		mappedQuestions[i] = models.Question{
			Name:        question.Name,
			Description: question.Description,
			Variants:    mapVariants(question.Variants),
		}
	}

	return mappedQuestions
}

func mapVariants(variants []CreateVariantInput) []models.Variant {
	mappedVariants := make([]models.Variant, len(variants))

	for i, variant := range variants {
		mappedVariants[i] = models.Variant{
			Name:      variant.Name,
			IsCorrect: variant.IsCorrect,
		}
	}

	return mappedVariants
}
func SetGetAllTests(tests []models.Test, count int64) *GetAllTestsResponse {
	return &GetAllTestsResponse{
		Tests: tests,
		Count: count,
	}
}

type GetAllTestsRequest struct {
	UserId uint `json:"user_id" validate:"required"`
	Limit  int  `json:"limit" validate:"required"`
	Offset int  `json:"offset" `
}

type GetTestByIdRequest struct {
	TestId uint `json:"test_id" validate:"required"`
}

type GetAllTestsResponse struct {
	Tests []models.Test `json:"tests"`
	Count int64         `json:"count"`
}

type GetTestByIdResponse struct {
	models.Test
	Role string `json:"user_role"`
}

type CreateTestRequest struct {
	Name      string                `json:"name" validate:"required"`
	Questions []CreateQuestionInput `json:"questions" validate:"required"`
}

type CreateAnonymusTestRequest struct {
	Name      string                `json:"name" validate:"required"`
	Questions []CreateQuestionInput `json:"questions"`
}

type CreateQuestionInput struct {
	Name        string               `json:"name" validate:"required"`
	Description string               `json:"description" validate:"required"`
	Variants    []CreateVariantInput `json:"variants"`
}

type CreateVariantInput struct {
	Name      string `json:"name" validate:"required"`
	IsCorrect bool   `json:"is_correct" validate:"required"`
}

type UpdateTestActiveStatus struct {
	TestId   uint `json:"test_id" validate:"required"`
	IsActive bool `json:"is_active"`
}
