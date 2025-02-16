package storage

import (
	"gorm.io/gorm"
)

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
