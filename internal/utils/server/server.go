package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/server/configs"
	delivery "github.com/server/delivery/http"
	"github.com/server/internal/utils/db/postgresql"
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
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: false,
	}

	corsHandler := cors.New(corsOptions).Handler(s.router)
	s.log.Sugar().Infof("Server run on http://localhost" + s.addr)
	return http.ListenAndServe(s.addr, corsHandler)
}

func (s *api) FillEndpoints() {
	delivery.NewAuth(s.router, s.log, s.db, s.cfg)
	delivery.NewTestManager(s.log, s.db, s.router)
	delivery.NewValidateResult(s.db, s.router, s.log)
	delivery.NewUser(s.log, s.db, s.router, s.cfg)
}
