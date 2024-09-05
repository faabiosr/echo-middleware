/*
 * Copyright (c) Fabio da Silva Ribeiro <faabiosr@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package middleware

import (
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// ZapLogConfig defines the config for Uber ZapLog middleware.
type ZapLogConfig struct {
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

	// Logger it is a zap logger
	Logger *zap.Logger

	// Skipper defines a function to skip middleware.
	Skipper mw.Skipper
}

// DefaultZapLogConfig is the default Uber ZapLog middleware config.
var DefaultZapLogConfig = ZapLogConfig{
	FieldMap: defaultFields,
	Logger: func() *zap.Logger {
		lg, _ := zap.NewProduction()
		return lg
	}(),
	Skipper: mw.DefaultSkipper,
}

// ZapLog returns a middleware that logs HTTP requests.
func ZapLog() echo.MiddlewareFunc {
	return ZapLogWithConfig(DefaultZapLogConfig)
}

// ZapLogWithConfig returns a Uber ZapLog middleware with config.
// See: `ZapLog()`.
func ZapLogWithConfig(cfg ZapLogConfig) echo.MiddlewareFunc {
	// Defaults
	if cfg.Skipper == nil {
		cfg.Skipper = DefaultZapLogConfig.Skipper
	}

	if cfg.Logger == nil {
		cfg.Logger = DefaultZapLogConfig.Logger
	}

	if len(cfg.FieldMap) == 0 {
		cfg.FieldMap = DefaultZapLogConfig.FieldMap
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) (err error) {
			if cfg.Skipper(ec) {
				return next(ec)
			}

			zFields := []zap.Field{}
			logFields, err := mapFields(ec, next, cfg.FieldMap)

			for k, v := range logFields {
				field := zap.Any(k, v)
				zFields = append(zFields, field)
			}

			cfg.Logger.With(zFields...).Info("handle request")

			return
		}
	}
}
