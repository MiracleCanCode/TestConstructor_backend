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

func (s *TestValidator) Validate(test *models.Test) (*float64, error) {
	exampleTest, err := s.testRepo.GetTestById(test.ID)
	if err != nil {
		s.logger.Error("Failed to get test by ID", zap.Error(err))
		return nil, errors.New("failed to fetch test")
	}

	err = s.testRepo.IncrementCountUserPast(test.ID, int(exampleTest.CountUserPast))
	if err != nil {
		s.logger.Error("Failed to increment count user past", zap.Error(err))
		return nil, errors.New("failed to increment count user past")
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
		s.logger.Warn("No answers provided")
		return nil, nil
	}

	percentage := (float64(totalCorrect) / float64(totalAnswers)) * 100
	return &percentage, nil
}
