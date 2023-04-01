package middleware

import (
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
)

// SLogConfig defines the config for golang Structured Log middleware.
type SLogConfig struct {
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

	// Logger it is a slog logger
	Logger *slog.Logger

	// Skipper defines a function to skip middleware.
	Skipper mw.Skipper
}

// DefaultSLogConfig is the default golang Structured Log middleware config.
var DefaultSLogConfig = SLogConfig{
	FieldMap: defaultFields,
	Logger:   slog.Default(),
	Skipper:  mw.DefaultSkipper,
}

// SLog returns a middleware that logs HTTP requests.
func SLog() echo.MiddlewareFunc {
	return SLogWithConfig(DefaultSLogConfig)
}

// SLogWithConfig returns a golang Structured Log middleware with config.
// See: `SLog()`.
func SLogWithConfig(cfg SLogConfig) echo.MiddlewareFunc {
	// Defaults
	if cfg.Skipper == nil {
		cfg.Skipper = DefaultSLogConfig.Skipper
	}

	if cfg.Logger == nil {
		cfg.Logger = DefaultSLogConfig.Logger
	}

	if len(cfg.FieldMap) == 0 {
		cfg.FieldMap = DefaultSLogConfig.FieldMap
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) (err error) {
			if cfg.Skipper(ec) {
				return next(ec)
			}

			cFields := []interface{}{}
			logFields, err := mapFields(ec, next, cfg.FieldMap)

			for k, v := range logFields {
				cFields = append(append(cFields, k), v)
			}

			cfg.Logger.With(cFields...).Info("handle request")

			return
		}
	}
}
