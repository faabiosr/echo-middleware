/*
 * Copyright (c) Fabio da Silva Ribeiro <faabiosr@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package middleware

import (
	charm "github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

// CharmLogConfig defines the config for CharmBracelet Log middleware.
type CharmLogConfig struct {
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
	// - @route
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

	// Logger it is a charm logger
	Logger *charm.Logger

	// Skipper defines a function to skip middleware.
	Skipper mw.Skipper
}

// DefaultCharmLogConfig is the default CharmBracelet Log middleware config.
var DefaultCharmLogConfig = CharmLogConfig{
	FieldMap: defaultFields,
	Logger:   charm.Default(),
	Skipper:  mw.DefaultSkipper,
}

// CharmLog returns a middleware that logs HTTP requests.
func CharmLog() echo.MiddlewareFunc {
	return CharmLogWithConfig(DefaultCharmLogConfig)
}

// CharmLogWithConfig returns a CharmBracelet Log middleware with config.
// See: `CharmLog()`.
func CharmLogWithConfig(cfg CharmLogConfig) echo.MiddlewareFunc {
	// Defaults
	if cfg.Skipper == nil {
		cfg.Skipper = DefaultCharmLogConfig.Skipper
	}

	if cfg.Logger == nil {
		cfg.Logger = DefaultCharmLogConfig.Logger
	}

	if len(cfg.FieldMap) == 0 {
		cfg.FieldMap = DefaultCharmLogConfig.FieldMap
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

			cfg.Logger.Info("handle request", cFields...)

			return
		}
	}
}

// CharmLogRecoverFn returns a CharmLog recover log function to print panic
// errors.
func CharmLogRecoverFn(logger *charm.Logger) mw.LogErrorFunc {
	return func(_ echo.Context, err error, stack []byte) error {
		logger.Error(
			"panic recover",
			"stacktrace", string(stack),
			"error", err,
		)

		return err
	}
}
