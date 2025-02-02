package main

import (
	"github.com/MiracleCanCode/example_configuration_logger"
	"github.com/server/configs"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/server"
	"go.uber.org/zap"
)

func main() {
	log := logger.Logger(logger.DefaultLoggerConfig())
	defer func() {
		if err := log.Sync(); err != nil {
			log.Error("Failed to sync logger", zap.Error(err))
		}
	}()
	conf, err := configs.Load(log)
	if err != nil {
		log.Error("Failed to load config", zap.Error(err))
		return
	}
	db, err := postgresql.New(conf, log)
	if err != nil {
		log.Error("Failed to initialize db", zap.Error(err))
		return
	}
	defer db.Close()
	app := server.New(db, log, conf)

	if err := app.RunApp(); err != nil {
		log.Error("Failed to run server", zap.Error(err))
		return
	}

}
