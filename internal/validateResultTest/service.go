package validateresulttest

import (
	"github.com/server/internal/getTest"
	"github.com/server/models"
	"github.com/server/pkg/db/postgresql"
	"go.uber.org/zap"
)

type Service struct {
	db             *postgresql.Db
	getTestService *getTest.Service
	logger         *zap.Logger
}

func NewService(db *postgresql.Db, logger *zap.Logger, service *getTest.Service) *Service {
	return &Service{
		db:             db,
		logger:         logger,
		getTestService: service,
	}
}

func (s *Service) Validate(test *models.Test) (*float64, error) {
	exampleTest, err := s.getTestService.GetById(test.ID)
	if err != nil {
		s.logger.Error("Failed to get test by id, error: " + err.Error())
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
