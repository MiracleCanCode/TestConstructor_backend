package usecases

import (
	"fmt"

	"github.com/server/entity"
	"go.uber.org/zap"
)

type TestManagerRepoV2Interface interface {
	GetTestById(id uint) (*entity.Test, error)
	IncrementCountUserPast(testId uint, count int) error
}

type TestValidator struct {
	testManagerRepo TestManagerRepoV2Interface
	logger          *zap.Logger
}

func NewTestValidator(
	testManagerRepo TestManagerRepoV2Interface,
	logger *zap.Logger,
) *TestValidator {
	return &TestValidator{
		testManagerRepo: testManagerRepo,
		logger:          logger,
	}
}

func (s *TestValidator) Validate(test *entity.Test) (*float64, error) {
	exampleTest, err := s.testManagerRepo.GetTestById(test.ID)
	if err != nil {
		return nil, fmt.Errorf("Validate: failed to get test by ID: %w", err)
	}

	err = s.testManagerRepo.IncrementCountUserPast(test.ID, int(exampleTest.CountUserPast))
	if err != nil {
		return nil, fmt.Errorf("Validate: failed to increment count user past: %w", err)
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
