/*
 * Copyright (c) Fabio da Silva Ribeiro <faabiosr@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package middleware

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
)

type ctxkey struct{ name string }

// reqIDKey key used to store the request-id in context.
var reqIDKey = &ctxkey{"request-id"}

// RequestIDConfig alias for emw.RequestIDConfig
type RequestIDConfig = emw.RequestIDConfig

// DefaultRequestIDConfig is the default RequestID middleware config, based on
// the echo.RequestIDConfig but with uuid generator instead.
var DefaultRequestIDConfig = RequestIDConfig{
	Skipper:          emw.DefaultSkipper,
	Generator:        uuidGen,
	RequestIDHandler: requestIDHandler,
	TargetHeader:     echo.HeaderXRequestID,
}

// RequestID returns a middleware that reads or generates a new request id and
// returns to response, also stores in context.
func RequestID() echo.MiddlewareFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestIDWithConfig uses the echo.RequestIDWithConfig under the hood with
// custom generator and sets the request id in context.
func RequestIDWithConfig(cfg RequestIDConfig) echo.MiddlewareFunc {
	// Defaults
	if cfg.Skipper == nil {
		cfg.Skipper = emw.DefaultRequestIDConfig.Skipper
	}

	if cfg.Generator == nil {
		cfg.Generator = uuidGen
	}

	if cfg.RequestIDHandler == nil {
		cfg.RequestIDHandler = requestIDHandler
	}

	if cfg.TargetHeader == "" {
		cfg.TargetHeader = echo.HeaderXRequestID
	}

	return emw.RequestIDWithConfig(cfg)
}

// uuidGen generates a random uuid.V4.
func uuidGen() string {
	return uuid.Must(uuid.NewV4()).String()
}

// requestIDHandler sets the received request-id into request context.
func requestIDHandler(ec echo.Context, rid string) {
	req := ec.Request()
	ctx := context.WithValue(req.Context(), reqIDKey, rid)

	ec.SetRequest(req.WithContext(ctx))
}

// RequestIDValue returns the value stored in the context, otherwise returns an
// empty string
func RequestIDValue(ctx context.Context) string {
	v, _ := ctx.Value(reqIDKey).(string)
	return v
}
