package postgresql

import (
	"fmt"
	"os"

	"time"

	"github.com/server/configs"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func New(conf *configs.Config, log *zap.Logger) (*Db, error) {
	maxRetries := 5
	initialDelay := 1 * time.Second
	var db *gorm.DB
	var err error

	debugMode, ok := os.LookupEnv("DEBUG_MODE")
	isDebug := false

	if ok && debugMode == "true" {
		isDebug = true
	}
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(conf.DB), &gorm.Config{})
		if err == nil {
			break
		}
		delay := initialDelay * time.Duration(1<<i)
		log.Error("Failed conn to db", zap.Error(err), zap.Int("retry", i), zap.Duration("delay", delay))
		time.Sleep(delay)
	}

	if err != nil {
		log.Error("Max retries reached, failed to connect to database.", zap.Error(err))
		return nil, fmt.Errorf("max retries reached: %w", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Error("Failed to create db instance", zap.Error(err))
		return nil, fmt.Errorf("failed to create db instance: %w", err)
	}

	if isDebug {
		db.Debug()
	}
	sqlDb.SetConnMaxIdleTime(30 * time.Minute)
	sqlDb.SetMaxOpenConns(50)
	sqlDb.SetConnMaxLifetime(30 * time.Minute)

	log.Info("DB connected!")
	return &Db{
		db,
	}, nil
}

func (d *Db) Close() {
	sqlDb, err := d.DB.DB()
	if err == nil {
		sqlDb.Close()
	}
}
