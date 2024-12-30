package createTest

import "github.com/server/models"

func MapCreateTestRequestToModel(req *CreateTestRequest) *models.Test {
	test := &models.Test{
		Name:        req.Name,
		AuthorLogin: req.AuthorLogin,
		Questions:   mapQuestions(req.Questions),
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

type CreateTestRequest struct {
	Name        string                `json:"name" validate:"required"`
	AuthorLogin string                `json:"author_login" validate:"required"`
	Questions   []CreateQuestionInput `json:"questions" validate:"required"`
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
