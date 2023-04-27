package middleware

import (
	"bytes"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestZeroLogWithConfig(t *testing.T) {
	ec := postCtx(t)
	b := new(bytes.Buffer)
	logger := log.Output(zerolog.ConsoleWriter{Out: b, NoColor: true})

	config := ZeroLogConfig{
		Logger:   logger,
		FieldMap: testFields,
	}

	_ = ZeroLogWithConfig(config)(testHandler)(ec)

	tests := []struct {
		str string
		err string
	}{
		{"handle request", "invalid log: handle request info not found"},
		{"id=123", "invalid log: request id not found"},
		{`remote_ip=http://foo.bar`, "invalid log: remote ip not found"},
		{`uri=http://some/foo/456?name=john`, "invalid log: uri not found"},
		{"host=some", "invalid log: host not found"},
		{"method=POST", "invalid log: method not found"},
		{"status=200", "invalid log: status not found"},
		{"latency=", "invalid log: latency not found"},
		{"latency_human=", "invalid log: latency_human not found"},
		{"bytes_in=0", "invalid log: bytes_in not found"},
		{"bytes_out=4", "invalid log: bytes_out not found"},
		{"path=/foo/456", "invalid log: path not found"},
		{`route=/foo/:id`, "invalid log: route not found"},
		{"protocol=HTTP/1.1", "invalid log: protocol not found"},
		{`referer=http://foo.bar`, "invalid log: referer not found"},
		{"user_agent=cli-agent", "invalid log: user_agent not found"},
		{"user=admin", "invalid log: header user not found"},
		{"filter_name=john", "invalid log: query filter_name not found"},
		{"username=doejohn", "invalid log: form field username not found"},
		{"session=A1B2C3", "invalid log: cookie session not found"},
	}

	for _, test := range tests {
		if !strings.Contains(b.String(), test.str) {
			t.Error(test.err)
		}
	}
}

func TestZeroLog(t *testing.T) {
	ec := reqCtx(t)
	_ = ZeroLog()(testHandler)(ec)
}

func TestZeroLogWithEmptyConfig(t *testing.T) {
	ec := reqCtx(t)
	_ = ZeroLogWithConfig(ZeroLogConfig{})(testHandler)(ec)
}

func TestZeroLogWithSkipper(t *testing.T) {
	ec := reqCtx(t)

	config := DefaultZeroLogConfig
	config.Skipper = func(c echo.Context) bool {
		return true
	}

	_ = ZeroLogWithConfig(config)(testHandler)(ec)
}

func TestZeroLogRetrievesAnError(t *testing.T) {
	ec := errCtx(t)
	b := new(bytes.Buffer)
	logger := log.Output(zerolog.ConsoleWriter{Out: b, NoColor: true})

	config := ZeroLogConfig{
		Logger: logger,
	}

	_ = ZeroLogWithConfig(config)(testHandler)(ec)

	res := b.String()

	if !strings.Contains(res, "status=500") {
		t.Errorf("invalid log: wrong status code")
	}

	if !strings.Contains(res, `error=error`) {
		t.Errorf("invalid log: error not found")
	}
}
