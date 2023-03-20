package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"views-servive/api"
	"views-servive/config"
	"views-servive/services"
)

type App struct {
	server   *http.Server
	mainCtx  context.Context
	settings config.Settings
}

func NewApp(mainCtx context.Context, settings config.Settings) *App {
	return &App{
		mainCtx:  mainCtx,
		settings: settings,
	}
}

func (a *App) Start() {
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			log.Printf("Server didn't start: %v\n", err)
		}
	}()
}

func (a *App) Stop(getContext func(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)) error {
	serverCtx, cancelServerCtx := getContext(a.mainCtx, time.Second*15)
	defer cancelServerCtx()

	err := a.server.Shutdown(serverCtx)
	if err != nil {
		log.Printf("Server didn't stop: %v", err)
		return err
	}

	return nil
}

func (a *App) InitServices() error {
	parseService := services.NewParseServiceImp(a.settings)
	viewsService := services.NewViewsServiceImp(parseService, a.settings)

	a.server = api.NewServer(
		a.mainCtx,
		a.settings,
		parseService,
		viewsService,
	)

	return nil
}
