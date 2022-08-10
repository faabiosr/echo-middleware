package middleware

import (
	"testing"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
)

func TestRequestIDWithConfig(t *testing.T) {
	ec := reqCtx(t)

	cfg := emw.RequestIDConfig{}
	_ = RequestIDWithConfig(cfg)(testHandler)(ec)

	rid := RequestIDValue(ec.Request().Context())
	header := ec.Response().Header().Get(echo.HeaderXRequestID)

	if rid != header {
		t.Error("id does not match")
	}
}

func TestRequestID(t *testing.T) {
	ec := reqCtx(t)

	_ = RequestID()(testHandler)(ec)

	rid := RequestIDValue(ec.Request().Context())
	header := ec.Response().Header().Get(echo.HeaderXRequestID)

	if rid != header {
		t.Error("id does not match")
	}
}
