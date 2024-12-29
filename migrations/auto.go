package main

import (
	"github.com/MiracleCanCode/example_configuration_logger/pkg/logger"
	"github.com/MiracleCanCode/zaperr"
	"github.com/server/configs"
	"github.com/server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log := logger.Logger(logger.DefaultLoggerConfig())
	handleErrors := zaperr.NewZaperr(log)
	db, err := gorm.Open(postgres.Open(configs.LoadConfig(log, handleErrors).DB), &gorm.Config{})

	handleErrors.LogError(err, "Failed to open db")

	db.AutoMigrate(&models.User{}, &models.Test{}, &models.Question{}, &models.Variant{})
}
