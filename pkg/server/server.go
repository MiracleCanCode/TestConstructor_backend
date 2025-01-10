package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/server/configs"
	"github.com/server/internal/auth"
	testmanager "github.com/server/internal/testmanager"
	"github.com/server/internal/user"
	validateresulttest "github.com/server/internal/validateResultTest"
	"github.com/server/pkg/db/postgresql"
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
	auth.New(s.router, s.log, s.db, s.cfg)
	testmanager.New(s.log, s.db, s.router)
	validateresulttest.New(s.db, s.router, s.log)
	user.New(s.log, s.db, s.router, s.cfg)
}
