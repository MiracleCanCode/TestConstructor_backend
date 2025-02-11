package main

import (
	"os"

	"github.com/server/configs"
	"github.com/server/entity"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	log := logger.GetInstance()
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

	connPostgres := db.Connection()

	if err := connPostgres.AutoMigrate(&entity.User{}, &entity.Test{}, &entity.Question{}, &entity.Variant{}); err != nil {
		log.Error("Error migration", zap.Error(err))
		os.Exit(1)
	}

	log.Info("Migrations completed")

}
