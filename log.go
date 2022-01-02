package middleware

import (
	"net/url"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// Log middlewares constants.
const (
	logID           = "@id"
	logRemoteIP     = "@remote_ip"
	logURI          = "@uri"
	logHost         = "@host"
	logMethod       = "@method"
	logPath         = "@path"
	logProtocol     = "@protocol"
	logReferer      = "@referer"
	logUserAgent    = "@user_agent"
	logStatus       = "@status"
	logError        = "@error"
	logLatency      = "@latency"
	logLatencyHuman = "@latency_human"
	logBytesIn      = "@bytes_in"
	logBytesOut     = "@bytes_out"
	logHeader       = "@header"
	logQuery        = "@query"
	logForm         = "@form"
	logCookie       = "@cookie"
)

// int conversion base.
const base = 10

// mapTags maps the field tags with its respective data.
func mapTags(ctx echo.Context, latency time.Duration) map[string]interface{} {
	fields := map[string]interface{}{}

	req := ctx.Request()
	res := ctx.Response()

	id := req.Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = res.Header().Get(echo.HeaderXRequestID)
	}

	fields[logID] = id
	fields[logRemoteIP] = ctx.RealIP()
	fields[logURI] = req.RequestURI
	fields[logHost] = req.Host
	fields[logMethod] = req.Method

	path := req.URL.Path
	if path == "" {
		path = "/"
	}

	fields[logPath] = path
	fields[logProtocol] = req.Proto
	fields[logReferer] = req.Referer()
	fields[logUserAgent] = req.UserAgent()
	fields[logStatus] = res.Status

	cl := req.Header.Get(echo.HeaderContentLength)

	if cl == "" {
		cl = "0"
	}

	fields[logLatency] = strconv.FormatInt(int64(latency), base)
	fields[logLatencyHuman] = latency.String()
	fields[logBytesIn] = cl
	fields[logBytesOut] = strconv.FormatInt(res.Size, base)

	fields = logParams(logQuery, ctx.QueryParams(), fields)

	if params, err := ctx.FormParams(); err == nil {
		fields = logParams(logForm, params, fields)
	}

	return fields
}

func logParams(prefix string, values url.Values, fields map[string]interface{}) map[string]interface{} {
	for k := range values {
		key := prefix + ":" + k
		fields[key] = values.Get(k)
	}

	return fields
}
