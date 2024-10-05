package app

import (
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// DefaultPort is the default port the server listens on.
const DefaultPort = 8080

// New creates a new Echo instance.
func New(_ *zap.Logger) *echo.Echo {
	app := echo.New()
	app.HideBanner = true
	app.HidePort = true

	app.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	return app
}

// GetPort returns the port the server is listening on.
func GetPort(logger *zap.Logger) int64 {
	rawPort := os.Getenv("PORT")
	if rawPort == "" {
		return DefaultPort
	}

	port, err := strconv.ParseInt(rawPort, 10, 64)
	if err != nil {
		logger.Warn("Failed to parse port", zap.String("raw_port", rawPort))
		logger.Warn("Using default port", zap.Int64("port", DefaultPort))

		return DefaultPort
	}

	return port
}
