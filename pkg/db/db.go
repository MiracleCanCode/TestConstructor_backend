package db

import (
	"context"

	"gorm.io/gorm"
)

type DBInterface interface {
	Connection() *gorm.DB
	Close()
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
	RollbackTransaction(tx *gorm.DB) error
	CommitTransaction(tx *gorm.DB) error
}

type DB struct {
	*gorm.DB
}

func (d *DB) Connection() (*gorm.DB, error) {
	return d.DB, nil
}

func (d *DB) Close() {
	sqlDb, err := d.DB.DB()
	if err == nil {
		sqlDb.Close()
	}
}
