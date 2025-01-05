package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/server/configs"
	"github.com/server/internal/auth"
	"github.com/server/internal/createTest"
	"github.com/server/internal/getTest"
	validateresulttest "github.com/server/internal/validateResultTest"
	"github.com/server/pkg/db/postgresql"
	"github.com/server/pkg/middleware"
	"go.uber.org/zap"
)

type api struct {
	addr         string
	router       *mux.Router
	db           *postgresql.Db
	log          *zap.Logger
	cfg          *configs.Config
}

func New(db *postgresql.Db, logger *zap.Logger, cfg *configs.Config) *api {
	router := mux.NewRouter()
	router.Use(middleware.CORS)
	return &api{
		addr:   cfg.PORT,
		router: router,
		db:     db,
		log:    logger,
	}
}

func (s *api) RunApp() error {
	s.log.Sugar().Infof("Server run on http://localhost" + s.addr)
	return http.ListenAndServe(s.addr, s.router)
}

func (s *api) FillEndpoints() {
	auth.New(s.router, s.log, s.db, s.cfg)
	createTest.New(s.log, s.db, s.router)
	getTest.New(s.log, s.db, s.router)
	validateresulttest.New(s.db, s.router, s.log)
}

func (s *api) GenerateDocs() {

}
