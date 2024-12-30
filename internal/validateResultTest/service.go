package validateresulttest

import (
	"github.com/server/internal/getTest"
	"github.com/server/models"
	"github.com/server/pkg/db"
	"go.uber.org/zap"
)

type ValidateResultTestService struct {
	db             *db.Db
	getTestService *getTest.GetTestService
	logger         *zap.Logger
}

func NewValidateResultTestService(db *db.Db, logger *zap.Logger, service *getTest.GetTestService) *ValidateResultTestService {
	return &ValidateResultTestService{
		db:             db,
		logger:         logger,
		getTestService: service,
	}
}

func (s *ValidateResultTestService) Validate(test *models.Test) (*float64, error) {
	exampleTest, err := s.getTestService.GetTestById(test.ID)
	if err != nil {
		s.logger.Error("Failed to get test by id, error: " + err.Error())
		return nil, err
	}

	var totalCorrect int
	var totalAnswers int

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
