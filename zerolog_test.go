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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestZeroLogWithConfig(t *testing.T) {
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

	logger := log.Output(zerolog.ConsoleWriter{Out: b, NoColor: true})

	fields := DefaultZeroLogConfig.FieldMap
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

	config := ZeroLogConfig{
		Logger:   logger,
		FieldMap: fields,
	}

	ZeroLogWithConfig(config)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)

	res := b.String()

	tests := []struct {
		str string
		err string
	}{
		{"handle request", "invalid log: handle request info not found"},
		{"id=123", "invalid log: request id not found"},
		{`remote_ip=http://foo.bar`, "invalid log: remote ip not found"},
		{`uri=http://some?name=john`, "invalid log: uri not found"},
		{"host=some", "invalid log: host not found"},
		{"method=POST", "invalid log: method not found"},
		{"status=200", "invalid log: status not found"},
		{"latency=", "invalid log: latency not found"},
		{"latency_human=", "invalid log: latency_human not found"},
		{"bytes_in=0", "invalid log: bytes_in not found"},
		{"bytes_out=4", "invalid log: bytes_out not found"},
		{"path=/", "invalid log: path not found"},
		{"protocol=HTTP/1.1", "invalid log: protocol not found"},
		{`referer=http://foo.bar`, "invalid log: referer not found"},
		{"user_agent=cli-agent", "invalid log: user_agent not found"},
		{"user=admin", "invalid log: header user not found"},
		{"filter_name=john", "invalid log: query filter_name not found"},
		{"username=doejohn", "invalid log: form field username not found"},
		{"session=A1B2C3", "invalid log: cookie session not found"},
	}

	for _, test := range tests {
		if !strings.Contains(res, test.str) {
			t.Error(test.err)
		}
	}
}

func TestZeroLog(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ZeroLog()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)
}

func TestZeroLogWithEmptyConfig(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ZeroLogWithConfig(ZeroLogConfig{})(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)
}

func TestZeroLogWithSkipper(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := DefaultZeroLogConfig
	config.Skipper = func(c echo.Context) bool {
		return true
	}

	ZeroLogWithConfig(config)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)
}

func TestZeroLogRetrievesAnError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	b := new(bytes.Buffer)

	logger := log.Output(zerolog.ConsoleWriter{Out: b, NoColor: true})

	config := ZeroLogConfig{
		Logger: logger,
	}

	ZeroLogWithConfig(config)(func(c echo.Context) error {
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
