package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nicklasfrahm-dev/appkit/logging"
	"go.uber.org/zap"
	"kycelis.dev/core/pkg/app"
)

// version is injected by the build process.
var version = "dev"

const shutdownTimeout = 5 * time.Second

func main() {
	logger := logging.NewLogger()

	logger.Info("Starting application", zap.String("version", version))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGPIPE,
	)

	errorChannel := make(chan error, 1)

	server := app.New(logger)

	go startServer(logger, server, errorChannel)

	select {
	case err := <-errorChannel:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Failed to run HTTP server", zap.Error(err))
		}

		logger.Warn("Terminating application due to server shutdown")
	case sig := <-signals:
		logger.Info("Received signal", zap.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("Failed to shutdown server gracefully", zap.Error(err))
		}
	}
}

func startServer(logger *zap.Logger, srv *echo.Echo, errorChannel chan<- error) {
	port := app.GetPort(logger)

	logger.Info("Starting server", zap.Int64("port", port))

	errorChannel <- srv.Start(fmt.Sprintf(":%d", port))
}
