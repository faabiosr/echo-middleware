package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"go.opencensus.io/stats/view"
)

func TestOpenCensus(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	OpenCensus()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)
}

func TestOpenCensusWithEmptyConfig(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	OpenCensusWithConfig(OpenCensusConfig{})(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)
}

func TestOpenCensusWithSkipper(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := DefaultOpenCensusConfig
	config.Skipper = func(c echo.Context) bool {
		return true
	}

	OpenCensusWithConfig(config)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)
}

func TestOpenCensusWithWrongView(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic when register views")
		}
	}()

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/some", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := OpenCensusConfig{
		Views: []*view.View{{}},
	}

	OpenCensusWithConfig(config)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})(c)
}
