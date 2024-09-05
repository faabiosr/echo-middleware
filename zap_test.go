/*
 * Copyright (c) Fabio da Silva Ribeiro <faabiosr@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapLogWithConfig(t *testing.T) {
	ec := postCtx(t)
	logger, logs := observer.New(zap.InfoLevel)

	config := ZapLogConfig{
		Logger:   zap.New(logger),
		FieldMap: testFields,
	}

	_ = ZapLogWithConfig(config)(testHandler)(ec)

	entry := logs.All()[0]
	ectx := entry.ContextMap()

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"msg", entry.Message, "handle request"},
		{"id", ectx["id"], "123"},
		{"remote_ip", ectx["remote_ip"], "http://foo.bar"},
		{"uri", ectx["uri"], "http://some/foo/456?name=john"},
		{"host", ectx["host"], "some"},
		{"method", ectx["method"], "POST"},
		{"status", ectx["status"], int64(http.StatusOK)},
		{"bytes_in", ectx["bytes_in"], "0"},
		{"bytes_out", ectx["bytes_out"], "4"},
		{"path", ectx["path"], "/foo/456"},
		{"route", ectx["route"], "/foo/:id"},
		{"protocol", ectx["protocol"], "HTTP/1.1"},
		{"referer", ectx["referer"], "http://foo.bar"},
		{"user_agent", ectx["user_agent"], "cli-agent"},
		{"user", ectx["user"], "admin"},
		{"filter_name", ectx["filter_name"], "john"},
		{"username", ectx["username"], "doejohn"},
		{"session", ectx["session"], "A1B2C3"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("entry_%s", tt.name), func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("expect '%s' as '%v', got '%v'", tt.name, tt.want, tt.got)
			}
		})
	}
}

func TestZapLog(t *testing.T) {
	ec := reqCtx(t)
	_ = ZapLog()(testHandler)(ec)
}

func TestZapLogWithEmptyConfig(t *testing.T) {
	ec := reqCtx(t)
	_ = ZapLogWithConfig(ZapLogConfig{})(testHandler)(ec)
}

func TestZapLogWithSkipper(t *testing.T) {
	ec := reqCtx(t)

	config := DefaultZapLogConfig
	config.Skipper = func(echo.Context) bool {
		return true
	}

	_ = ZapLogWithConfig(config)(testHandler)(ec)
}

func TestZapLogRetrievesAnError(t *testing.T) {
	ec := errCtx(t)
	logger, logs := observer.New(zap.InfoLevel)

	config := ZapLogConfig{
		Logger: zap.New(logger),
	}

	_ = ZapLogWithConfig(config)(testHandler)(ec)

	entry := logs.All()[0]
	ectx := entry.ContextMap()

	if ectx["status"] != int64(http.StatusInternalServerError) {
		t.Errorf("invalid log: wrong status code")
	}

	if _, ok := ectx["error"]; !ok {
		t.Errorf("invalid log: error not found")
	}
}

func TestZapLogRecoverFn(t *testing.T) {
	ec := panicCtx(t)
	obsLog, logs := observer.New(zap.InfoLevel)
	logger := zap.New(obsLog)

	rec := emw.RecoverWithConfig(emw.RecoverConfig{
		LogErrorFunc: ZapLogRecoverFn(logger),
	})

	config := ZapLogConfig{
		Logger: logger,
	}

	_ = ZapLogWithConfig(config)(rec(testHandler))(ec)

	entry := logs.All()[0]
	ectx := entry.ContextMap()

	if _, ok := ectx["error"]; !ok {
		t.Errorf("invalid log: error not found")
	}

	entry = logs.All()[1]
	ectx = entry.ContextMap()

	if ectx["status"] != int64(http.StatusInternalServerError) {
		t.Errorf("invalid log: wrong status code")
	}
}
