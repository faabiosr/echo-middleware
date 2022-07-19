package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapLogWithConfig(t *testing.T) {
	e := echo.New()

	form := url.Values{}
	form.Add("username", "doejohn")

	req := httptest.NewRequest(echo.POST, "http://some?name=john", strings.NewReader(form.Encode()))

	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Add("Referer", "http://foo.bar")
	req.Header.Add("User-Agent", "cli-agent")
	req.Header.Add(echo.HeaderXForwardedFor, "http://foo.bar")
	req.Header.Add("user", "admin")
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: "A1B2C3",
	})

	rec := httptest.NewRecorder()
	rec.Header().Add(echo.HeaderXRequestID, "123")

	ec := e.NewContext(req, rec)

	fields := DefaultZapLogConfig.FieldMap
	fields["empty"] = ""
	fields["id"] = logID
	fields["path"] = logPath
	fields["protocol"] = logProtocol
	fields["referer"] = logReferer
	fields["user_agent"] = logUserAgent
	fields["store"] = logHeaderPrefix + "store"
	fields["filter_name"] = logQueryPrefix + "name"
	fields["username"] = logFormPrefix + "username"
	fields["session"] = logCookiePrefix + "session"
	fields["bytes_in"] = logBytesIn
	fields["bytes_out"] = logBytesOut
	fields["referer"] = logReferer
	fields["user"] = logHeaderPrefix + "user"

	logger, logs := observer.New(zap.InfoLevel)

	config := ZapLogConfig{
		Logger:   zap.New(logger),
		FieldMap: fields,
	}

	_ = ZapLogWithConfig(config)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(ec)

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
		{"uri", ectx["uri"], "http://some?name=john"},
		{"host", ectx["host"], "some"},
		{"method", ectx["method"], "POST"},
		{"status", ectx["status"], int64(http.StatusOK)},
		{"bytes_in", ectx["bytes_in"], "0"},
		{"bytes_out", ectx["bytes_out"], "4"},
		{"path", ectx["path"], "/"},
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
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	ec := e.NewContext(req, rec)

	_ = ZapLog()(func(ec echo.Context) error {
		return ec.String(http.StatusOK, "test")
	})(ec)
}

func TestZapLogWithEmptyConfig(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	ec := e.NewContext(req, rec)

	_ = ZapLogWithConfig(ZapLogConfig{})(func(ec echo.Context) error {
		return ec.String(http.StatusOK, "test")
	})(ec)
}

func TestZapLogWithSkipper(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := DefaultZapLogConfig
	config.Skipper = func(c echo.Context) bool {
		return true
	}

	_ = ZapLogWithConfig(config)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)
}

func TestZapLogRetrievesAnError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	ec := e.NewContext(req, rec)

	logger, logs := observer.New(zap.InfoLevel)

	config := ZapLogConfig{
		Logger: zap.New(logger),
	}

	_ = ZapLogWithConfig(config)(func(ec echo.Context) error {
		return errors.New("error")
	})(ec)

	entry := logs.All()[0]
	ectx := entry.ContextMap()

	if ectx["status"] != int64(http.StatusInternalServerError) {
		t.Errorf("invalid log: wrong status code")
	}

	if _, ok := ectx["error"]; !ok {
		t.Errorf("invalid log: error not found")
	}
}
