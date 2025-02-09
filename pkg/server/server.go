package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/server/configs"
	delivery "github.com/server/domain/http"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/middleware"
	"go.uber.org/zap"
)

type api struct {
	addr   string
	router *mux.Router
	db     *postgresql.Db
	log    *zap.Logger
	cfg    *configs.Config
}

func New(db *postgresql.Db, logger *zap.Logger, cfg *configs.Config) *api {
	router := mux.NewRouter()

	return &api{
		addr:   cfg.PORT,
		router: router,
		db:     db,
		log:    logger,
	}
}

func (s *api) RunApp() error {
	s.FillEndpoints()
	handler := middleware.DefaultCORSMiddleware()(s.router)
	s.log.Info("Server started")
	return http.ListenAndServe(s.addr, handler)
}

func (s *api) FillEndpoints() {
	delivery.NewAuthHandler(s.router, s.log, s.db, s.cfg)
	delivery.NewTestManagerHandler(s.log, s.db, s.router)
	delivery.NewValidateResultHandler(s.db, s.router, s.log)
	delivery.NewUserHandler(s.log, s.db, s.router, s.cfg)
	s.router.Handle("/metrics", promhttp.Handler())
}
