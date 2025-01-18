package postgresql

import (
	"github.com/server/configs"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func New(conf *configs.Config, log *zap.Logger) *Db {
	db, err := gorm.Open(postgres.Open(conf.DB), &gorm.Config{})

	if err != nil {
		log.Error("Failed conn to db", zap.Error(err))
	}

	return &Db{
		db,
	}
}
