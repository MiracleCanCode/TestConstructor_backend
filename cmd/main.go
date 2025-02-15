package main

import (
	"github.com/server/configs"
	"github.com/server/internal/transport"
	"github.com/server/pkg/logger"

	"github.com/server/pkg/storage/postgresql"
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

	app := transport.New(connPostgres, log, conf)

	if err := app.RunApp(); err != nil {
		log.Error("Failed to run server", zap.Error(err))
		return
	}

}
