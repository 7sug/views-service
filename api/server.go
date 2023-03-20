package api

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"views-servive/api/handlers"
	"views-servive/config"
	"views-servive/services"
)

type Server struct {
	server   *http.Server
	settings *config.Settings
}

func NewServer(ctx context.Context, settings config.Settings, parseService services.ParseServiceImp, viewsService services.ViewsServiceImp) *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/ping", handlers.PingHandler()).Methods(http.MethodGet)
	router.HandleFunc("/test-parse", handlers.TestParseHandler(parseService)).Methods(http.MethodGet)
	router.HandleFunc("/views", handlers.ViewsHandler(viewsService)).Methods(http.MethodPost)

	return &http.Server{
		Addr: fmt.Sprintf(":%d", settings.Port),
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		Handler: router,
	}
}
