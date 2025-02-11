package main

import (
	"github.com/server/configs"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/logger"
	"github.com/server/pkg/server"
	"go.uber.org/zap"
)

func main() {
	log := logger.GetInstance()

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
	connPostgres := db.Connection()

	app := server.New(connPostgres, log, conf)

	if err := app.RunApp(); err != nil {
		log.Error("Failed to run server", zap.Error(err))
		return
	}

}
