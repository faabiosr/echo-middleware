package middleware

import (
	"strconv"
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
		"remote_ip": "@remote_ip",
		"uri":       "@uri",
		"host":      "@host",
		"method":    "@method",
		"status":    "@status",
		"latency":   "@latency",
		"error":     "@error",
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

			req := ctx.Request()
			res := ctx.Response()
			start := time.Now()

			if err = next(ctx); err != nil {
				ctx.Error(err)
			}

			stop := time.Now()
			entry := cfg.Logger

			for k, v := range cfg.FieldMap {
				if v == "" {
					continue
				}

				switch v {
				case "@id":
					id := req.Header.Get(echo.HeaderXRequestID)

					if id == "" {
						id = res.Header().Get(echo.HeaderXRequestID)
					}

					entry = entry.WithField(k, id)
				case "@remote_ip":
					entry = entry.WithField(k, ctx.RealIP())
				case "@uri":
					entry = entry.WithField(k, req.RequestURI)
				case "@host":
					entry = entry.WithField(k, req.Host)
				case "@method":
					entry = entry.WithField(k, req.Method)
				case "@path":
					p := req.URL.Path

					if p == "" {
						p = "/"
					}

					entry = entry.WithField(k, p)
				case "@protocol":
					entry = entry.WithField(k, req.Proto)
				case "@referer":
					entry = entry.WithField(k, req.Referer())
				case "@user_agent":
					entry = entry.WithField(k, req.UserAgent())
				case "@status":
					entry = entry.WithField(k, res.Status)
				case "@error":
					if err != nil {
						entry = entry.WithField(k, err)
					}
				case "@latency":
					l := stop.Sub(start)
					entry = entry.WithField(k, strconv.FormatInt(int64(l), 10))
				case "@latency_human":
					entry = entry.WithField(k, stop.Sub(start).String())
				case "@bytes_in":
					cl := req.Header.Get(echo.HeaderContentLength)

					if cl == "" {
						cl = "0"
					}

					entry = entry.WithField(k, cl)
				case "@bytes_out":
					entry = entry.WithField(k, strconv.FormatInt(res.Size, 10))
				default:
					switch {
					case strings.HasPrefix(v, "@header:"):
						entry = entry.WithField(k, ctx.Request().Header.Get(v[8:]))
					case strings.HasPrefix(v, "@query:"):
						entry = entry.WithField(k, ctx.QueryParam(v[7:]))
					case strings.HasPrefix(v, "@form:"):
						entry = entry.WithField(k, ctx.FormValue(v[6:]))
					case strings.HasPrefix(v, "@cookie:"):
						cookie, err := ctx.Cookie(v[8:])
						if err == nil {
							entry = entry.WithField(k, cookie.Value)
						}
					}
				}
			}

			entry.Print("handle request")

			return
		}
	}
}
