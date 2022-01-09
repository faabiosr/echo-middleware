package middleware

import (
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ZeroLogConfig defines the config for ZeroLog middleware.
type ZeroLogConfig struct {
	// FieldMap set a list of fields with tags
	//
	// Tags to constructed the logger fields.
	//
	// - @id (Request ID)
	// - @remote_ip
	// - @uri
	// - @host
	// - @method
	// - @path
	// - @protocol
	// - @referer
	// - @user_agent
	// - @status
	// - @error
	// - @latency (In nanoseconds)
	// - @latency_human (Human readable)
	// - @bytes_in (Bytes received)
	// - @bytes_out (Bytes sent)
	// - @header:<NAME>
	// - @query:<NAME>
	// - @form:<NAME>
	// - @cookie:<NAME>
	FieldMap map[string]string

	// Logger it is a zerolog logger
	Logger zerolog.Logger

	// Skipper defines a function to skip middleware.
	Skipper mw.Skipper
}

// DefaultZeroLogConfig is the default ZeroLog middleware config.
var DefaultZeroLogConfig = ZeroLogConfig{
	FieldMap: map[string]string{
		"remote_ip": logRemoteIP,
		"uri":       logURI,
		"host":      logHost,
		"method":    logMethod,
		"status":    logStatus,
		"latency":   logLatency,
		"error":     logError,
	},
	Logger:  log.Logger,
	Skipper: mw.DefaultSkipper,
}

// ZeroLog returns a middleware that logs HTTP requests.
func ZeroLog() echo.MiddlewareFunc {
	return ZeroLogWithConfig(DefaultZeroLogConfig)
}

// ZeroLogWithConfig returns a ZeroLog middleware with config.
// See: `ZeroLog()`.
func ZeroLogWithConfig(cfg ZeroLogConfig) echo.MiddlewareFunc {
	// Defaults
	if cfg.Skipper == nil {
		cfg.Skipper = DefaultZeroLogConfig.Skipper
	}

	if len(cfg.FieldMap) == 0 {
		cfg.FieldMap = DefaultZeroLogConfig.FieldMap
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			if cfg.Skipper(ctx) {
				return next(ctx)
			}

			logFields, err := mapFields(ctx, next, cfg.FieldMap)

			cfg.Logger.Info().
				Fields(logFields).
				Msg("handle request")

			return
		}
	}
}
