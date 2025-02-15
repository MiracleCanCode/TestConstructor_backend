package postgresql

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/server/configs"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBInterface interface {
	Connection() *gorm.DB
	Close()
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
	RollbackTransaction(tx *gorm.DB) error
	CommitTransaction(tx *gorm.DB) error
}

type PostgresDB struct {
	db *gorm.DB
}

var (
	instance *PostgresDB
	once     sync.Once
)

func New(conf *configs.Config, log *zap.Logger) (DBInterface, error) {
	var err error
	once.Do(func() {
		maxRetries := 5
		initialDelay := 1 * time.Second
		var db *gorm.DB

		debugMode, ok := os.LookupEnv("DEBUG_MODE")
		isDebug := ok && debugMode == "true"

		for i := 0; i < maxRetries; i++ {
			db, err = gorm.Open(postgres.Open(conf.DB), &gorm.Config{})
			if err == nil {
				break
			}
			delay := initialDelay * time.Duration(1<<i)
			log.Error("Failed to connect to DB", zap.Error(err), zap.Int("retry", i), zap.Duration("delay", delay))
			time.Sleep(delay)
		}

		if err != nil {
			log.Error("Max retries reached, failed to connect to database.", zap.Error(err))
			return
		}

		sqlDb, err := db.DB()
		if err != nil {
			log.Error("Failed to create db instance", zap.Error(err))
			return
		}

		if isDebug {
			db = db.Debug()
		}

		sqlDb.SetConnMaxIdleTime(30 * time.Minute)
		sqlDb.SetMaxOpenConns(150)
		sqlDb.SetConnMaxLifetime(30 * time.Minute)

		instance = &PostgresDB{db: db}
	})

	if err != nil {
		return nil, fmt.Errorf("New postgres: %w", err)
	}
	return instance, nil
}

func (p *PostgresDB) Connection() *gorm.DB {
	return p.db
}

func (p *PostgresDB) Close() {
	sqlDb, err := p.db.DB()
	if err == nil {
		sqlDb.Close()
	}
}

func (p *PostgresDB) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx := p.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("Begin transaction: %w", tx.Error)
	}
	return tx, nil
}

func (p *PostgresDB) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (p *PostgresDB) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}
