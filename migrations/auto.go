package main

import (
	"github.com/MiracleCanCode/example_configuration_logger/pkg/logger"
	"github.com/server/configs"
	"github.com/server/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log := logger.Logger(logger.DefaultLoggerConfig())
	db, err := gorm.Open(postgres.Open(configs.Load(log).DB), &gorm.Config{})

	if err != nil {
		log.Error("Failed to open db", zap.Error(err))
		return
	}

	if err := db.AutoMigrate(&models.User{}, &models.Test{}, &models.Question{}, &models.Variant{}); err != nil {
		log.Error("Error migration", zap.Error(err))
		return
	}
}
