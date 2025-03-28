package dtos

import "github.com/server/entity"

type GetTestResponse struct {
	ID            uint
	Name          string                `json:"name"`
	AuthorLogin   string                `json:"author_login"`
	UserID        uint                  `json:"user_id"`
	IsActive      bool                  `json:"is_active"`
	CountUserPast uint                  `json:"count_user_past"`
	Questions     []GetQuestionResponse `json:"questions"`
	Role          string                `json:"user_role"`
}

type GetQuestionResponse struct {
	ID          uint
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Variants    []GetVariantResponse `json:"variants"`
}

type GetVariantResponse struct {
	ID   uint
	Name string `json:"name"`
}

func MapTestToGetTestResponse(test *entity.Test, role string, userID uint) *GetTestResponse {
	questions := make([]GetQuestionResponse, len(test.Questions))

	for i, question := range test.Questions {
		variants := make([]GetVariantResponse, len(question.Variants))
		for j, variant := range question.Variants {
			variants[j] = GetVariantResponse{
				ID:   variant.ID,
				Name: variant.Name,
			}
		}

		questions[i] = GetQuestionResponse{
			ID:          question.ID,
			Name:        question.Name,
			Description: question.Description,
			Variants:    variants,
		}
	}

	return &GetTestResponse{
		ID:            test.ID,
		Name:          test.Name,
		AuthorLogin:   test.AuthorLogin,
		UserID:        userID,
		IsActive:      test.IsActive,
		CountUserPast: test.CountUserPast,
		Questions:     questions,
		Role:          role,
	}
}

func MapCreateTestRequestToModel(req *CreateTestRequest, userId uint) entity.Test {
	test := entity.Test{
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
	*entity.Test
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
