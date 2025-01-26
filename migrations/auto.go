package main

import (
	"os"

	"github.com/MiracleCanCode/example_configuration_logger/pkg/logger"
	"github.com/server/configs"
	"github.com/server/internal/models"
	"github.com/server/pkg/db/postgresql"
	"go.uber.org/zap"
)

func main() {
	log := logger.Logger(logger.DefaultLoggerConfig())
	cfg, err := configs.Load(log)
	if err != nil {
		log.Error("Failed to load config", zap.Error(err))
		os.Exit(1)
	}

	db, err := postgresql.New(cfg, log)
	if err != nil {
		log.Error("Failed to open db", zap.Error(err))
		os.Exit(1)
	}

	defer db.Close()

	if err := db.AutoMigrate(&models.User{}, &models.Test{}, &models.Question{}, &models.Variant{}); err != nil {
		log.Error("Error migration", zap.Error(err))
		os.Exit(1)
	}

	log.Info("Migrations completed")

}
