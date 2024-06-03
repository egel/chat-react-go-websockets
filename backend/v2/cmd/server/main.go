package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"chatserver/internal/app"
	"chatserver/pkg/logging"

	"github.com/rs/zerolog/log"
)

func main() {
	// add logging
	output := logging.ConfigConsoleOutput()
	log.Logger = log.Output(output).With().Caller().Logger()

	const (
		SERVER_HOST = "localhost"
		SERVER_PORT = 8000
	)

	appOptions := app.NewAppOptions(1) // 0 = No GracefulShutdownTime
	serverAddress := fmt.Sprintf("%s:%d", SERVER_HOST, SERVER_PORT)
	app := app.Initialize(serverAddress, appOptions)

	go func() {
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("ListenAndServe ERROR")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of X.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), appOptions.GracefulShutdownTime)
	defer cancel()
	if err := app.Server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server Shutdown")
	}
	// catching ctx.Done(). timeout of X seconds.
	select {
	case <-ctx.Done():
		log.Info().
			Dur("GracefulShutdownTime", app.Options.GracefulShutdownTime).
			Msg("timeout of before closing server")
	}
	log.Info().Msg("Server exiting")
}
