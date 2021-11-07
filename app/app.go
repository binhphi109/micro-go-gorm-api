package app

import (
	"github.com/sirupsen/logrus"

	"github.com/sample/sample-server/core"
	"github.com/sample/sample-server/store"
)

type IApp interface {
	Company() CompanyApp
	User() UserApp
}

type App struct {
	company *CompanyApp
	user    *UserApp
	Store   *store.Store
	Logger  *logrus.Logger
	Config  *core.Config
}

func New(store *store.Store, logger *logrus.Logger, config *core.Config) *App {
	return &App{
		Store:  store,
		Logger: logger,
		Config: config,
	}
}
