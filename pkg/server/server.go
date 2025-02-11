package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/server/configs"
	delivery "github.com/server/delivery/http"
	"github.com/server/delivery/http/middleware"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type api struct {
	addr   string
	router *mux.Router
	db     *gorm.DB
	log    *zap.Logger
	cfg    *configs.Config
}

func New(db *gorm.DB, logger *zap.Logger, cfg *configs.Config) *api {
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
