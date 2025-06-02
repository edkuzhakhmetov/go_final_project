package application

import (
	"context"
	"fmt"
	"net/http"

	"go1f/pkg/config"

	"github.com/edkuzhakhmetov/go_final_project/internal/api"
	"github.com/edkuzhakhmetov/go_final_project/internal/storage"
	"github.com/edkuzhakhmetov/go_final_project/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Application struct {
	cfg     *config.Config
	log     logrus.FieldLogger
	storage *storage.Storage
	handler *api.Handler
	router  *api.Router
}

func New() *Application {
	return &Application{}
}

func (app *Application) Run(ctx context.Context) error {
	var err error

	app.log = logger.GetLogger()
	app.log.Info("Starting application")

	app.cfg, err = config.LoadConfig()
	if err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}
	app.log.Info("Config loaded")

	app.storage, err = storage.NewStorage(app.cfg.ToDoDBFile)
	if err != nil {
		return fmt.Errorf("can't connect to db: %w", err)
	}
	app.log.Info("DB connected")

	err = app.storage.Migrate()
	if err != nil {
		return fmt.Errorf("can't migrate db: %w", err)
	}
	app.log.Info("DB migrated")

	app.handler = api.NewHandler(app.log, app.storage)
	app.log.Info("Handler created")

	app.router = api.NewRouter(app.handler)
	app.log.Info("Router created")

	app.log.Infof("Starting server on port %s", app.cfg.ToDoPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", app.cfg.ToDoPort), app.router.Mux)

	if err != nil {
		app.log.Errorf("error starting server: %v", err)
		return fmt.Errorf("error starting server: %w", err)
	}

	return nil
}
