package usecases

import (
	"errors"

	"github.com/server/internal/models"
	"github.com/server/internal/repository"
	"go.uber.org/zap"
)

type TestValidatorInterface interface {
	Validate(test *models.Test) (*float64, error)
}

type TestValidator struct {
	testRepo repository.TestManagerInterface
	logger   *zap.Logger
}

func NewTestValidator(
	testRepo repository.TestManagerInterface,
	logger *zap.Logger,
) *TestValidator {
	return &TestValidator{
		testRepo: testRepo,
		logger:   logger,
	}
}

func (tv *TestValidator) Validate(test *models.Test) (*float64, error) {
	exampleTest, err := tv.testRepo.GetTestById(test.ID)
	if err != nil {
		tv.logger.Error("Failed to get test by ID", zap.Error(err))
		return nil, errors.New("failed to fetch test")
	}

	var (
		totalCorrect int
		totalAnswers int
	)

	for _, question := range exampleTest.Questions {
		for _, userQuestion := range test.Questions {
			if question.ID != userQuestion.ID {
				continue
			}

			for _, variant := range question.Variants {
				for _, userVariant := range userQuestion.Variants {
					if variant.Name == userVariant.Name {
						totalAnswers++
						if variant.IsCorrect == userVariant.IsCorrect {
							totalCorrect++
						}
					}
				}
			}
		}
	}

	if totalAnswers == 0 {
		tv.logger.Warn("No answers provided")
		return nil, nil
	}

	percentage := (float64(totalCorrect) / float64(totalAnswers)) * 100
	return &percentage, nil
}
