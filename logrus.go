package middleware

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// LogrusConfig defines the config for Logrus middleware.
type LogrusConfig struct {
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

	// Logger it is a logrus logger
	Logger logrus.FieldLogger

	// Skipper defines a function to skip middleware.
	Skipper mw.Skipper
}

// DefaultLogrusConfig is the default Logrus middleware config.
var DefaultLogrusConfig = LogrusConfig{
	FieldMap: map[string]string{
		"remote_ip": logRemoteIP,
		"uri":       logURI,
		"host":      logHost,
		"method":    logMethod,
		"status":    logStatus,
		"latency":   logLatency,
		"error":     logError,
	},
	Logger:  logrus.StandardLogger(),
	Skipper: mw.DefaultSkipper,
}

// Logrus returns a middleware that logs HTTP requests.
func Logrus() echo.MiddlewareFunc {
	return LogrusWithConfig(DefaultLogrusConfig)
}

// LogrusWithConfig returns a Logrus middleware with config.
// See: `Logrus()`.
func LogrusWithConfig(cfg LogrusConfig) echo.MiddlewareFunc {
	// Defaults
	if cfg.Skipper == nil {
		cfg.Skipper = DefaultLogrusConfig.Skipper
	}

	if cfg.Logger == nil {
		cfg.Logger = DefaultLogrusConfig.Logger
	}

	if len(cfg.FieldMap) == 0 {
		cfg.FieldMap = DefaultLogrusConfig.FieldMap
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			if cfg.Skipper(ctx) {
				return next(ctx)
			}

			start := time.Now()

			if err = next(ctx); err != nil {
				ctx.Error(err)
			}

			stop := time.Now()
			latency := stop.Sub(start)

			entry := cfg.Logger
			tags := mapTags(ctx, latency)

			for k, tag := range cfg.FieldMap {
				if tag == "" {
					continue
				}

				if value, ok := tags[tag]; ok {
					entry = entry.WithField(k, value)
					continue
				}

				if tag == logError && err != nil {
					entry = entry.WithField(k, err)
					continue
				}

				if strings.HasPrefix(tag, logHeader+":") {
					entry = entry.WithField(k, ctx.Request().Header.Get(tag[8:]))
					continue
				}

				if strings.HasPrefix(tag, logCookie+":") {
					cookie, err := ctx.Cookie(tag[8:])
					if err == nil {
						entry = entry.WithField(k, cookie.Value)
					}
				}
			}

			entry.Print("handle request")

			return
		}
	}
}
