// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/graze/golang-service/log"
	"github.com/stretchr/testify/assert"
)

func TestStructuredLogging(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()

	userRequest := newRequest("GET", "http://example.com")
	userRequest.Header.Add("X-Forwarded-For", "192.168.100.5")

	userAgentRequest := newRequest("GET", "http://example.com")
	userAgentRequest.Header.Add("User-Agent", "some user agent")

	referrerRequest := newRequest("GET", "http://example.com/test")
	referrerRequest.Header.Add("Referer", "http://google.com")

	cases := map[string]struct {
		request   *http.Request
		timestamp time.Time
		duration  time.Duration
		size      int
		message   string
		fields    map[string]interface{}
	}{
		"basic": {
			newRequest("GET", "http://example.com"),
			now,
			getDuration(t, "0.302s"),
			100,
			"GET / HTTP/1.1",
			map[string]interface{}{
				"module":          "request.handler",
				"tag":             "request_handled",
				"http.method":     "GET",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/",
				"http.path":       "/",
				"http.host":       "example.com",
				"http.status":     200,
				"http.bytes":      100,
				"dur":             0.302,
				"http.time":       now.Format(time.RFC3339Nano),
				"http.ref":        "",
				"http.user":       "",
				"http.user-agent": "",
				"transaction":     "test-123",
			},
		},
		"post path": {
			newRequest("POST", "http://example.com/path/here"),
			now,
			getDuration(t, "0.102s"),
			200,
			"POST /path/here HTTP/1.1",
			map[string]interface{}{
				"module":          "request.handler",
				"tag":             "request_handled",
				"http.method":     "POST",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/path/here",
				"http.path":       "/path/here",
				"http.host":       "example.com",
				"http.status":     200,
				"http.bytes":      200,
				"dur":             0.102,
				"http.time":       now.Format(time.RFC3339Nano),
				"http.ref":        "",
				"http.user":       "",
				"http.user-agent": "",
				"transaction":     "test-123",
			},
		},
		"strips params off method": {
			newRequest("GET", "http://example.com/token/1/test?apid=1&thing=2"),
			now,
			getDuration(t, "0.927321s"),
			300,
			"GET /token/1/test?apid=1&thing=2 HTTP/1.1",
			map[string]interface{}{
				"module":          "request.handler",
				"tag":             "request_handled",
				"http.method":     "GET",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/token/1/test?apid=1&thing=2",
				"http.path":       "/token/1/test",
				"http.host":       "example.com",
				"http.status":     200,
				"http.bytes":      300,
				"dur":             0.927321,
				"http.time":       now.Format(time.RFC3339Nano),
				"http.ref":        "",
				"http.user":       "",
				"http.user-agent": "",
				"transaction":     "test-123",
			},
		},
		"connect http2 test": {
			&http.Request{
				Method:     "CONNECT",
				Proto:      "HTTP/2.0",
				ProtoMajor: 2,
				ProtoMinor: 0,
				URL:        &url.URL{Host: "www.example.com:443"},
				Host:       "www.example.com:443",
				RemoteAddr: "192.168.100.5:9843",
			},
			now,
			getDuration(t, "0.927s"),
			400,
			"CONNECT www.example.com:443 HTTP/2.0",
			map[string]interface{}{
				"module":          "request.handler",
				"tag":             "request_handled",
				"http.method":     "CONNECT",
				"http.protocol":   "HTTP/2.0",
				"http.uri":        "www.example.com:443",
				"http.path":       "www.example.com:443",
				"http.host":       "www.example.com:443",
				"http.status":     200,
				"http.bytes":      400,
				"dur":             0.927,
				"http.time":       now.Format(time.RFC3339Nano),
				"http.ref":        "",
				"http.user":       "192.168.100.5",
				"http.user-agent": "",
			},
		},
		"handles referrer": {
			referrerRequest,
			now,
			getDuration(t, "0.019s"),
			500,
			"GET /test HTTP/1.1",
			map[string]interface{}{
				"module":          "request.handler",
				"tag":             "request_handled",
				"http.method":     "GET",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/test",
				"http.path":       "/test",
				"http.host":       "example.com",
				"http.status":     200,
				"http.bytes":      500,
				"dur":             0.019,
				"http.time":       now.Format(time.RFC3339Nano),
				"http.ref":        "http://google.com",
				"http.user":       "",
				"http.user-agent": "",
				"transaction":     "test-123",
			},
		},
		"handles x-forwarded-for": {
			userRequest,
			now,
			getDuration(t, "0.019s"),
			600,
			"GET / HTTP/1.1",
			map[string]interface{}{
				"module":          "request.handler",
				"tag":             "request_handled",
				"http.method":     "GET",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/",
				"http.path":       "/",
				"http.host":       "example.com",
				"http.status":     200,
				"http.bytes":      600,
				"dur":             0.019,
				"http.time":       now.Format(time.RFC3339Nano),
				"http.ref":        "",
				"http.user":       "192.168.100.5",
				"http.user-agent": "",
				"transaction":     "test-123",
			},
		},
		"handles User-Agent": {
			userAgentRequest,
			now,
			getDuration(t, "0.019s"),
			600,
			"GET / HTTP/1.1",
			map[string]interface{}{
				"module":          "request.handler",
				"tag":             "request_handled",
				"http.method":     "GET",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/",
				"http.path":       "/",
				"http.host":       "example.com",
				"http.status":     200,
				"http.bytes":      600,
				"dur":             0.019,
				"http.time":       now.Format(time.RFC3339Nano),
				"http.ref":        "",
				"http.user":       "",
				"http.user-agent": "some user agent",
				"transaction":     "test-123",
			},
		},
	}

	logger := log.New().With(log.KV{"transaction": "test-123"})
	base, ok := logger.(*log.LoggerEntry)
	if !ok {
		t.Error("Unable to cast logger to log.LoggerEntry")
	}
	hook := test.NewLocal(base.Logger)
	local := logger.With(log.KV{"module": "request.handler"})

	for k, tc := range cases {
		hook.Reset()
		rec := httptest.NewRecorder()
		responseLogger := &responseLogger{w: rec}
		writeStructuredLog(responseLogger, local, tc.request, *tc.request.URL, tc.timestamp, tc.duration, http.StatusOK, tc.size)
		assert.Equal(t, 1, len(hook.Entries), "test %s - Has Log Entry", k)
		assert.Equal(t, log.InfoLevel, hook.LastEntry().Level, "test %s - Has Log Level", k)
		assert.Equal(t, tc.message, hook.LastEntry().Message, "test %s - Has Message", k)
		for f, v := range tc.fields {
			assert.Contains(t, hook.LastEntry().Data, f, "test %s - Has Field: %s", k, f)
			if _, ok := v.(float64); ok {
				assert.InDelta(t, v, hook.LastEntry().Data[f], 0.0001, "test %s - Field: %s", k, f)
			} else {
				assert.Equal(t, v, hook.LastEntry().Data[f], "test %s - Field: %s", k, f)
			}
		}
	}
}
