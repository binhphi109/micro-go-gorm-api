package core

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/sample/sample-server/store"
)

// App is the App object that runs HTTP server.
type Server struct {
	Config *Config
	Db     gorm.DB
	Logger *logrus.Logger
	Router *chi.Mux
	Store  *store.Store
}

// New is the constructor for App struct.
func NewServer() *Server {
	router := chi.NewRouter()
	logger := NewLogger()
	config := NewConfig()

	db, err := ConnectDatabase(config)
	if err != nil {
		logrus.WithError(err).Fatal("error connecting database")
		return nil
	}

	s := store.New(db, *logger)

	server := &Server{
		Router: router,
		Config: config,
		Logger: logger,
		Db:     *db,
		Store:  s,
	}

	return server
}

func (s *Server) Start() {
	server := &http.Server{
		Handler: s.Router,
	}

	s.Logger.Infof("Starting HTTP server")

	err := server.ListenAndServe()
	if err != nil {
		s.Logger.WithError(err).Fatal("error starting HTTP server")
	}
}
