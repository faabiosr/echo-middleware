package middleware_test

import (
	"os"

	charm "github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"

	middleware "github.com/faabiosr/echo-middleware"
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

// This example registers the ZapLog middleware with default configuration.
func ExampleZapLog() {
	e := echo.New()

	// Middleware
	e.Use(middleware.ZapLog())
}

// This example registers the ZapLog middleware with custom configuration.
func ExampleZapLogWithConfig() {
	e := echo.New()

	// Custom ZapLog logger instance
	logger, _ := zap.NewProduction()

	// Middleware
	logConfig := middleware.ZapLogConfig{
		Logger: logger,
		FieldMap: map[string]string{
			"uri":    "@uri",
			"host":   "@host",
			"method": "@method",
			"status": "@status",
		},
	}

	e.Use(middleware.ZapLogWithConfig(logConfig))
}

// This example registers the CharmBracelet Log middleware with default configuration.
func ExampleCharmLog() {
	e := echo.New()

	// Middleware
	e.Use(middleware.CharmLog())
}

// This example registers the CharmBracelet Log middleware with custom configuration.
func ExampleCharmLogWithConfig() {
	e := echo.New()

	// Middleware
	logConfig := middleware.CharmLogConfig{
		Logger: charm.Default(),
		FieldMap: map[string]string{
			"uri":    "@uri",
			"host":   "@host",
			"method": "@method",
			"status": "@status",
		},
	}

	e.Use(middleware.CharmLogWithConfig(logConfig))
}

// This example registers the RequestID middleware with default configuration.
func ExampleRequestID() {
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestID())
}

// This example registers the RequestID middleware with custom configuration.
func ExampleRequestIDWithConfig() {
	e := echo.New()

	// Middleware
	config := middleware.RequestIDConfig{
		TargetHeader: echo.HeaderXRequestID,
	}

	e.Use(middleware.RequestIDWithConfig(config))
}
