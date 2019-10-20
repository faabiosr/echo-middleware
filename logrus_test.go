package middleware

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func TestLogrusWithConfig(t *testing.T) {
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

	c := e.NewContext(req, rec)
	b := new(bytes.Buffer)

	logger := logrus.StandardLogger()
	logger.Out = b

	fields := DefaultLogrusConfig.FieldMap
	fields["empty"] = ""
	fields["id"] = "@id"
	fields["path"] = "@path"
	fields["protocol"] = "@protocol"
	fields["referer"] = "@referer"
	fields["user_agent"] = "@user_agent"
	fields["store"] = "@header:store"
	fields["filter_name"] = "@query:name"
	fields["username"] = "@form:username"
	fields["session"] = "@cookie:session"
	fields["latency_human"] = "@latency_human"
	fields["bytes_in"] = "@bytes_in"
	fields["bytes_out"] = "@bytes_out"
	fields["referer"] = "@referer"
	fields["user"] = "@header:user"

	config := LogrusConfig{
		Logger:   logger,
		FieldMap: fields,
	}

	LogrusWithConfig(config)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)

	res := b.String()

	if !strings.Contains(res, "handle request") {
		t.Error("invalid log: handle request info not found")
	}

	if !strings.Contains(res, "id=123") {
		t.Error("invalid log: request id not found")
	}

	if !strings.Contains(res, `remote_ip="http://foo.bar"`) {
		t.Error("invalid log: remote ip not found")
	}

	if !strings.Contains(res, `uri="http://some?name=john"`) {
		t.Error("invalid log: uri not found")
	}

	if !strings.Contains(res, "host=some") {
		t.Error("invalid log: host not found")
	}

	if !strings.Contains(res, "method=POST") {
		t.Error("invalid log: method not found")
	}

	if !strings.Contains(res, "status=200") {
		t.Error("invalid log: status not found")
	}

	if !strings.Contains(res, "latency=") {
		t.Error("invalid log: latency not found")
	}

	if !strings.Contains(res, "latency_human=") {
		t.Error("invalid log: latency_human not found")
	}

	if !strings.Contains(res, "bytes_in=0") {
		t.Error("invalid log: bytes_in not found")
	}

	if !strings.Contains(res, "bytes_out=4") {
		t.Error("invalid log: bytes_out not found")
	}

	if !strings.Contains(res, "path=/") {
		t.Error("invalid log: path not found")
	}

	if !strings.Contains(res, "protocol=HTTP/1.1") {
		t.Error("invalid log: protocol not found")
	}

	if !strings.Contains(res, `referer="http://foo.bar"`) {
		t.Error("invalid log: referer not found")
	}

	if !strings.Contains(res, "user_agent=cli-agent") {
		t.Error("invalid log: user_agent not found")
	}

	if !strings.Contains(res, "user=admin") {
		t.Error("invalid log: header user not found")
	}

	if !strings.Contains(res, "filter_name=john") {
		t.Error("invalid log: query filter_name not found")
	}

	if !strings.Contains(res, "username=doejohn") {
		t.Error("invalid log: form field username not found")
	}

	if !strings.Contains(res, "session=A1B2C3") {
		t.Error("invalid log: cookie session not found")
	}
}

func TestLogrus(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	Logrus()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)
}

func TestLogrusWithEmptyConfig(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	LogrusWithConfig(LogrusConfig{})(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)
}

func TestLogrusWithSkipper(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := DefaultLogrusConfig
	config.Skipper = func(c echo.Context) bool {
		return true
	}

	LogrusWithConfig(config)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)
}

func TestLogrusRetrievesAnError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	b := new(bytes.Buffer)

	logger := logrus.StandardLogger()
	logger.Out = b

	config := LogrusConfig{
		Logger: logger,
	}

	LogrusWithConfig(config)(func(c echo.Context) error {
		return errors.New("error")
	})(c)

	res := b.String()

	if !strings.Contains(res, "status=500") {
		t.Errorf("invalid log: wrong status code")
	}

	if !strings.Contains(res, `error=error`) {
		t.Errorf("invalid log: error not found")
	}
}
