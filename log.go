package middleware

import (
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
	logHeaderPrefix = "@header:"
	logQueryPrefix  = "@query:"
	logFormPrefix   = "@form:"
	logCookiePrefix = "@cookie:"
)

// string to int base conversion.
const base = 10

// mapTags maps the log tags with its related data. Populate previously the
// key/value avoids the cyclomatic complexity of the log middlewares to
// identify each tag and value.
func mapTags(ec echo.Context, latency time.Duration) map[string]interface{} {
	tags := map[string]interface{}{}

	req := ec.Request()
	res := ec.Response()

	id := req.Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = res.Header().Get(echo.HeaderXRequestID)
	}

	tags[logID] = id
	tags[logRemoteIP] = ec.RealIP()
	tags[logURI] = req.RequestURI
	tags[logHost] = req.Host
	tags[logMethod] = req.Method

	path := req.URL.Path
	if path == "" {
		path = "/"
	}

	tags[logPath] = path
	tags[logProtocol] = req.Proto
	tags[logReferer] = req.Referer()
	tags[logUserAgent] = req.UserAgent()
	tags[logStatus] = res.Status
	tags[logLatency] = strconv.FormatInt(int64(latency), base)
	tags[logLatencyHuman] = latency.String()

	cl := req.Header.Get(echo.HeaderContentLength)
	if cl == "" {
		cl = "0"
	}

	tags[logBytesIn] = cl
	tags[logBytesOut] = strconv.FormatInt(res.Size, base)

	return tags
}
