package app

import (
	v1 "chatserver/internal/controller/v1"
	"net/http"
)

// App give a structure of what application is consist of
type App struct {
	Server  *http.Server
	Options *AppOptions
}

// NewApp creates a new instance of the application where it needs dependencies:
// http.Handler, serverAddress
func newApp(handler http.Handler, serverAddress string, options *AppOptions) *App {
	return &App{
		Server: &http.Server{
			Addr:    serverAddress,
			Handler: handler,
		},
		Options: options,
	}
}

func Initialize(serverAddress string, appoptions *AppOptions) *App {
	httpHandler := v1.NewHttpRouter()

	return newApp(httpHandler, serverAddress, appoptions)
}
