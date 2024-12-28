package main

import (
	"github.com/MiracleCanCode/zaperr"
	"github.com/server/configs"
	"github.com/server/models"
	"github.com/server/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log := logger.Logger()
	handleErrors := zaperr.NewZaperr(log)
	db, err := gorm.Open(postgres.Open(configs.LoadConfig(log, handleErrors).DB), &gorm.Config{})

	handleErrors.LogError(err, "Failed to open db")

	db.AutoMigrate(&models.User{}, &models.Test{}, &models.Question{})
}
