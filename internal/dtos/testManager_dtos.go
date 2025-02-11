package dtos

import "github.com/server/entity"

func MapTestModelToGetTestByIdResponse(req *entity.Test, userRole string, userId uint) *GetTestByIdResponse {
	return &GetTestByIdResponse{
		Test: entity.Test{
			Name:      req.Name,
			UserID:    userId,
			Questions: mapQuestionsGetTestById(req.Questions),
		},
		Role: userRole,
	}
}

func MapCreateTestRequestToModel(req *CreateTestRequest, userId uint) *entity.Test {
	test := &entity.Test{
		Name:      req.Name,
		UserID:    userId,
		Questions: mapQuestions(req.Questions),
	}

	return test
}

func mapQuestions(questions []CreateQuestionInput) []entity.Question {
	mappedQuestions := make([]entity.Question, len(questions))

	for i, question := range questions {
		mappedQuestions[i] = entity.Question{
			Name:        question.Name,
			Description: question.Description,
			Variants:    mapVariants(question.Variants),
		}
	}

	return mappedQuestions
}

func mapVariants(variants []CreateVariantInput) []entity.Variant {
	mappedVariants := make([]entity.Variant, len(variants))

	for i, variant := range variants {
		mappedVariants[i] = entity.Variant{
			Name:      variant.Name,
			IsCorrect: variant.IsCorrect,
		}
	}

	return mappedVariants
}

func mapQuestionsGetTestById(questions []entity.Question) []entity.Question {
	mappedQuestions := make([]entity.Question, len(questions))

	for i, question := range questions {
		mappedQuestions[i] = entity.Question{
			Name:        question.Name,
			Description: question.Description,
			Variants:    mapVariantsGetTestById(question.Variants),
		}
	}

	return mappedQuestions
}
func mapVariantsGetTestById(variants []entity.Variant) []entity.Variant {
	mappedVariants := make([]entity.Variant, len(variants))

	for i, variant := range variants {
		mappedVariants[i] = entity.Variant{
			Name: variant.Name,
		}
	}

	return mappedVariants
}

func SetGetAllTests(tests []entity.Test, count int64) *GetAllTestsResponse {
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
	Tests []entity.Test `json:"tests"`
	Count int64         `json:"count"`
}

type GetTestByIdResponse struct {
	entity.Test
	Role string `json:"user_role"`
}

type CreateTestRequest struct {
	Name      string                `json:"name" validate:"required"`
	Questions []CreateQuestionInput `json:"questions" validate:"required"`
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
