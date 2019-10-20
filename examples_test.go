package middleware_test

import (
	"os"

	middleware "github.com/faabiosr/echo-middleware"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// This example registers the ZeroLog middleware with default configuration.
func ExampleZeroLog() {
	e := echo.New()

	// Middleware
	e.Use(middleware.ZeroLog())
}

// This examples registers the ZeroLog middleware with customer configuration.
func ExampleZeroLogWithConfig() {
	e := echo.New()

	// Custom zerolog logger instance
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Middleware
	logConfig := middleware.ZeroLogConfig{
		Logger: logger,
		FieldMap: map[string]string{
			"uri":    "@uri",
			"host":   "@host",
			"method": "@method",
			"status": "@status",
		},
	}

	e.Use(middleware.ZeroLogWithConfig(logConfig))
}
