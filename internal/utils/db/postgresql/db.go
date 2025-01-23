package postgresql

import (
	"time"

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
		return nil
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Error("Failed to create db instance", zap.Error(err))
		return nil
	}

	if !conf.PRODACTION {
		db.Debug()
	}
	sqlDb.SetConnMaxIdleTime(30 * time.Minute)
	sqlDb.SetMaxOpenConns(50)
	sqlDb.SetConnMaxLifetime(30 * time.Minute)

	log.Info("DB connected!")
	return &Db{
		db,
	}
}
