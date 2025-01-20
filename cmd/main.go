package main

import (
	"github.com/MiracleCanCode/example_configuration_logger/pkg/logger"
	"github.com/server/configs"
	"github.com/server/internal/utils/db/postgresql"
	"github.com/server/internal/utils/server"
	"go.uber.org/zap"
)

func main() {
	log := logger.Logger(logger.DefaultLoggerConfig())
	conf := configs.Load(log)
	db := postgresql.New(conf, log)
	defer db.CloseConnection()
	app := server.New(db, log, conf)
	app.FillEndpoints()

	if err := app.RunApp(); err != nil {
		log.Error("Failed to run server", zap.Error(err))
		return
	}

}
