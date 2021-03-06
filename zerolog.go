package middleware

import (
	"strconv"
	"strings"
	"time"

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
		"remote_ip": "@remote_ip",
		"uri":       "@uri",
		"host":      "@host",
		"method":    "@method",
		"status":    "@status",
		"latency":   "@latency",
		"error":     "@error",
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

			req := ctx.Request()
			res := ctx.Response()
			start := time.Now()

			if err = next(ctx); err != nil {
				ctx.Error(err)
			}

			stop := time.Now()
			entry := cfg.Logger.Info()

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

					entry = entry.Str(k, id)
				case "@remote_ip":
					entry = entry.Str(k, ctx.RealIP())
				case "@uri":
					entry = entry.Str(k, req.RequestURI)
				case "@host":
					entry = entry.Str(k, req.Host)
				case "@method":
					entry = entry.Str(k, req.Method)
				case "@path":
					p := req.URL.Path

					if p == "" {
						p = "/"
					}

					entry = entry.Str(k, p)
				case "@protocol":
					entry = entry.Str(k, req.Proto)
				case "@referer":
					entry = entry.Str(k, req.Referer())
				case "@user_agent":
					entry = entry.Str(k, req.UserAgent())
				case "@status":
					entry = entry.Int(k, res.Status)
				case "@error":
					if err != nil {
						entry = entry.Err(err)
					}
				case "@latency":
					l := stop.Sub(start)
					entry = entry.Str(k, strconv.FormatInt(int64(l), 10))
				case "@latency_human":
					entry = entry.Str(k, stop.Sub(start).String())
				case "@bytes_in":
					cl := req.Header.Get(echo.HeaderContentLength)

					if cl == "" {
						cl = "0"
					}

					entry = entry.Str(k, cl)
				case "@bytes_out":
					entry = entry.Str(k, strconv.FormatInt(res.Size, 10))
				default:
					switch {
					case strings.HasPrefix(v, "@header:"):
						entry = entry.Str(k, ctx.Request().Header.Get(v[8:]))
					case strings.HasPrefix(v, "@query:"):
						entry = entry.Str(k, ctx.QueryParam(v[7:]))
					case strings.HasPrefix(v, "@form:"):
						entry = entry.Str(k, ctx.FormValue(v[6:]))
					case strings.HasPrefix(v, "@cookie:"):
						cookie, err := ctx.Cookie(v[8:])
						if err == nil {
							entry = entry.Str(k, cookie.Value)
						}
					}
				}
			}

			entry.Msg("handle request")

			return
		}
	}
}
