package validateresulttest

import (
	"github.com/server/internal/getTest"
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
