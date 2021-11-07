package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"

	"github.com/sample/sample-server/app"
	"github.com/sample/sample-server/core"
	"github.com/sample/sample-server/utils"
)

const (
	STATUS   = "status"
	StatusOk = "OK"
)

type API struct {
	App    *app.App
	Router *chi.Mux
	Logger *logrus.Logger
	Config *core.Config
}

func Init(server *core.Server) *API {
	api := &API{
		App:    app.New(server.Store, server.Logger, server.Config),
		Router: server.Router,
		Logger: server.Logger,
		Config: server.Config,
	}

	api.Router.Use(middleware.RequestID)
	api.Router.Use(middleware.RealIP)
	api.Router.Use(core.Logger("router", server.Logger))
	api.Router.Use(middleware.Recoverer)
	api.Router.Use(middleware.Timeout(time.Minute))

	api.InitUser()

	return api
}

func ReturnStatusOK(w http.ResponseWriter) {
	m := make(map[string]string)
	m[STATUS] = StatusOk
	w.Write([]byte(utils.MapToJSON(m)))
}
