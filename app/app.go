package app

import (
	"github.com/caseyrwebb/go-jwt-auth/app/data"
	"github.com/caseyrwebb/go-jwt-auth/app/router"
	"github.com/caseyrwebb/go-jwt-auth/app/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type App struct {
	Router *mux.Router
	DB     data.GoDB
}

func New(logger zap.Logger, configs *utils.Configurations) *App {
	a := &App{
		Router: mux.NewRouter(),
		DB:     &data.DB{},
	}

	router.InitRoutes(a.Router, a.DB, logger, configs)

	return a
}
