/*
 * Copyright (c) Fabio da Silva Ribeiro <faabiosr@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	// OpenCensusRequestCount counts the HTTP requests started.
	OpenCensusRequestCount = &view.View{
		Name:        "request_count",
		Description: "Count of HTTP request started",
		Measure:     ochttp.ServerRequestCount,
		Aggregation: view.Count(),
	}

	// OpenCensusRequestCountByMethod counts the HTTP requests by method.
	OpenCensusRequestCountByMethod = &view.View{
		Name:        "request_count_by_method",
		Description: "Server request count by HTTP method",
		TagKeys:     []tag.Key{ochttp.Method},
		Measure:     ochttp.ServerRequestCount,
		Aggregation: view.Count(),
	}

	// OpenCensusRequestCountByPath counts the HTTP requests by path.
	OpenCensusRequestCountByPath = &view.View{
		Name:        "request_count_by_path",
		Description: "Server request count by HTTP path",
		TagKeys:     []tag.Key{ochttp.Path},
		Measure:     ochttp.ServerRequestCount,
		Aggregation: view.Count(),
	}

	// OpenCensusResponseCountByStatusCode counts the HTTP requests by status code.
	OpenCensusResponseCountByStatusCode = &view.View{
		Name:        "response_count_by_status_code",
		Description: "Server response count by status code",
		TagKeys:     []tag.Key{ochttp.StatusCode},
		Measure:     ochttp.ServerLatency,
		Aggregation: view.Count(),
	}
)

// OpenCensusConfig defines the config for OpenCensus middleware.
type OpenCensusConfig struct {
	// View it is a OpenCensus Views list.
	Views []*view.View

	// Skipper defines a function to skip middleware.
	Skipper mw.Skipper
}

// DefaultOpenCensusConfig is the default OpenCensus middleware config.
var DefaultOpenCensusConfig = OpenCensusConfig{
	Views: []*view.View{
		OpenCensusRequestCount,
		OpenCensusRequestCountByMethod,
		OpenCensusRequestCountByPath,
		OpenCensusResponseCountByStatusCode,
	},
	Skipper: mw.DefaultSkipper,
}

// OpenCensus returns a middleware that collect HTTP requests and response
// metrics.
func OpenCensus() echo.MiddlewareFunc {
	return OpenCensusWithConfig(DefaultOpenCensusConfig)
}

// OpenCensusWithConfig returns a OpenCensus middleware with config.
// See: `OpenCensus()`.
func OpenCensusWithConfig(cfg OpenCensusConfig) echo.MiddlewareFunc {
	// Defaults
	if cfg.Skipper == nil {
		cfg.Skipper = DefaultOpenCensusConfig.Skipper
	}

	if len(cfg.Views) == 0 {
		cfg.Views = DefaultOpenCensusConfig.Views
	}

	if err := view.Register(cfg.Views...); err != nil {
		panic("echo: opencensus middleware register views failed")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			if cfg.Skipper(ctx) {
				return next(ctx)
			}

			handler := &ochttp.Handler{
				Handler: http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						ctx.SetRequest(r)
						ctx.SetResponse(echo.NewResponse(w, ctx.Echo()))
						err = next(ctx)
					},
				),
			}

			handler.ServeHTTP(ctx.Response(), ctx.Request())

			return
		}
	}
}
