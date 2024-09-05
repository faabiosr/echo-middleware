package middleware

import (
	"testing"

	"github.com/labstack/echo/v4"
	"go.opencensus.io/stats/view"
)

func TestOpenCensus(t *testing.T) {
	ec := reqCtx(t)

	_ = OpenCensus()(testHandler)(ec)
}

func TestOpenCensusWithEmptyConfig(t *testing.T) {
	ec := reqCtx(t)

	_ = OpenCensusWithConfig(OpenCensusConfig{})(testHandler)(ec)
}

func TestOpenCensusWithSkipper(t *testing.T) {
	ec := reqCtx(t)

	config := DefaultOpenCensusConfig
	config.Skipper = func(echo.Context) bool {
		return true
	}

	_ = OpenCensusWithConfig(config)(testHandler)(ec)
}

func TestOpenCensusWithWrongView(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic when register views")
		}
	}()

	ec := reqCtx(t)

	config := OpenCensusConfig{
		Views: []*view.View{{}},
	}

	_ = OpenCensusWithConfig(config)(testHandler)(ec)
}
