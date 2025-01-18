package usecases

import (
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/models"
	"github.com/server/repository"
	"go.uber.org/zap"
)

type ValidateResult struct {
	db     *postgresql.Db
	repo   *repository.TestManager
	logger *zap.Logger
}

func NewValidateResult(db *postgresql.Db, logger *zap.Logger) *ValidateResult {
	return &ValidateResult{
		db:     db,
		logger: logger,
		repo:   repository.NewTestManager(db),
	}
}

func (s *ValidateResult) Validate(test *models.Test) (*float64, error) {
	exampleTest, err := s.repo.GetTestById(test.ID)
	if err != nil {
		s.logger.Error("Failed to get test by id", zap.Error(err))
		return nil, err
	}

	var (
		totalCorrect int
		totalAnswers int
	)

	for _, question := range exampleTest.Questions {
		for _, variant := range question.Variants {
			for _, userQuestion := range test.Questions {
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
		return nil, nil
	}

	percentage := (float64(totalCorrect) / float64(totalAnswers)) * 100
	return &percentage, nil
}
