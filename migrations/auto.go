package main

import (
	"github.com/MiracleCanCode/example_configuration_logger/pkg/logger"
	"github.com/server/configs"
	"github.com/server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log := logger.Logger(logger.DefaultLoggerConfig())
	db, err := gorm.Open(postgres.Open(configs.Load(log).DB), &gorm.Config{})

	if err != nil {
		log.Error("Failed to open db, error:" + err.Error())
		return
	}

	if err := db.AutoMigrate(&models.User{}, &models.Test{}, &models.Question{}, &models.Variant{}); err != nil {
		log.Error("Error migration, error:" + err.Error())
		return
	}
}
