package middleware_test

import (
	"os"

	middleware "github.com/faabiosr/echo-middleware"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/stats/view"
)

// This example registers the ZeroLog middleware with default configuration.
func ExampleZeroLog() {
	e := echo.New()

	// Middleware
	e.Use(middleware.ZeroLog())
}

// This example registers the ZeroLog middleware with custom configuration.
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

// This example registers the Logrus middleware with default configuration.
func ExampleLogrus() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logrus())
}

// This example registers the Logrus middleware with custom configuration.
func ExampleLogrusWithConfig() {
	e := echo.New()

	// Custom logrus logger instance
	logger := logrus.New()

	// Middleware
	logConfig := middleware.LogrusConfig{
		Logger: logger,
		FieldMap: map[string]string{
			"uri":    "@uri",
			"host":   "@host",
			"method": "@method",
			"status": "@status",
		},
	}

	e.Use(middleware.LogrusWithConfig(logConfig))
}

// This example registers the OpenCensus middleware with default configuration.
func ExampleOpenCensus() {
	e := echo.New()

	// Middleware
	e.Use(middleware.OpenCensus())
}

// This example registers the OpenCensus middleware with custom configuration.
func ExampleOpenCensusWithConfig() {
	e := echo.New()

	// Middleware
	cfg := middleware.OpenCensusConfig{
		Views: []*view.View{
			middleware.OpenCensusRequestCount,
		},
	}

	e.Use(middleware.OpenCensusWithConfig(cfg))
}
