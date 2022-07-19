package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

var testFields = map[string]string{
	"remote_ip":     logRemoteIP,
	"uri":           logURI,
	"host":          logHost,
	"method":        logMethod,
	"status":        logStatus,
	"latency":       logLatency,
	"error":         logError,
	"empty":         "",
	"id":            logID,
	"path":          logPath,
	"protocol":      logProtocol,
	"referer":       logReferer,
	"user_agent":    logUserAgent,
	"store":         logHeaderPrefix + "store",
	"filter_name":   logQueryPrefix + "name",
	"username":      logFormPrefix + "username",
	"session":       logCookiePrefix + "session",
	"latency_human": logLatencyHuman,
	"bytes_in":      logBytesIn,
	"bytes_out":     logBytesOut,
	"user":          logHeaderPrefix + "user",
}

func testHandler(ec echo.Context) error {
	if v := ec.QueryParam("err"); v != "" {
		return errors.New("error")
	}

	return ec.String(http.StatusOK, "test")
}

func testCtx(t *testing.T, target string) echo.Context {
	t.Helper()

	req := httptest.NewRequest(echo.GET, target, nil)
	rec := httptest.NewRecorder()
	e := echo.New()

	return e.NewContext(req, rec)
}

func reqCtx(t *testing.T) echo.Context {
	return testCtx(t, "/some")
}

func errCtx(t *testing.T) echo.Context {
	return testCtx(t, "/some?err=1")
}

func postCtx(t *testing.T) echo.Context {
	t.Helper()

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

	e := echo.New()
	return e.NewContext(req, rec)
}
