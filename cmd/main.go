package main

import (
	"github.com/MiracleCanCode/zaperr"
	"github.com/server/cmd/api"
	"github.com/server/configs"
	"github.com/server/pkg/db"
	"github.com/server/pkg/logger"
)

func main() {
	log := logger.Logger()
	handleErrors := zaperr.NewZaperr(log)
	conf := configs.LoadConfig(log, handleErrors)
	db := db.NewDb(conf, log)
	app := api.New(db, log, conf, handleErrors)
	app.FillEndpoints()

	handleErrors.LogError(app.RunApp(), "")
}
