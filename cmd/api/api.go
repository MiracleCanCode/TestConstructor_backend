package api

import (
	"net/http"

	"github.com/MiracleCanCode/zaperr"
	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/internal/auth"
	"github.com/server/internal/createTest"
	"github.com/server/internal/getTest"
	validateresulttest "github.com/server/internal/validateResultTest"
	"github.com/server/pkg/db"
	"github.com/server/pkg/middleware"
	"go.uber.org/zap"
)

type api struct {
	addr         string
	router       *mux.Router
	db           *db.Db
	log          *zap.Logger
	cfg          *configs.Config
	handleErrors *zaperr.Zaperr
}

func New(db *db.Db, logger *zap.Logger, cfg *configs.Config, handleErrors *zaperr.Zaperr) *api {
	router := mux.NewRouter()
	router.Use(middleware.CORSMiddleware)
	return &api{
		addr:         cfg.PORT,
		router:       router,
		db:           db,
		log:          logger,
		handleErrors: handleErrors,
	}
}

func (s *api) RunApp() error {
	s.log.Sugar().Infof("Server run on http://localhost" + s.addr)
	return http.ListenAndServe(s.addr, s.router)
}

func (s *api) FillEndpoints() {
	auth.NewAuthHandler(s.router, s.log, s.db, s.cfg, s.handleErrors)
	createTest.NewCreateTestHandler(s.log, s.db, s.router, s.handleErrors)
	getTest.NewGetTestHandler(s.log, s.db, s.router, s.handleErrors)
	validateresulttest.NewValidateTestHandler(s.db, s.router, s.log)
}
