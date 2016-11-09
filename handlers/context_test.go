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

	"github.com/graze/golang-service/log"
	"github.com/stretchr/testify/assert"
)

func TestContextUpdatesTheRequestContext(t *testing.T) {
	t.Parallel()

	userRequest := newRequest("GET", "http://example.com")
	userRequest.Header.Add("X-Forwarded-For", "192.168.100.5")

	userAgentRequest := newRequest("GET", "http://example.com")
	userAgentRequest.Header.Add("User-Agent", "some user agent")

	referrerRequest := newRequest("GET", "http://example.com/test")
	referrerRequest.Header.Add("Referer", "http://google.com")

	cases := map[string]struct {
		request  *http.Request
		expected map[string]interface{}
		regex    map[string]string
	}{
		"basic": {
			newRequest("GET", "http://example.com"),
			map[string]interface{}{
				"http.method":     "GET",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/",
				"http.path":       "/",
				"http.host":       "example.com",
				"http.user":       "",
				"http.ref":        "",
				"http.user-agent": "",
			},
			map[string]string{
				"transaction": `(?:[0-9a-z]+-){4}[0-9a-z]+`,
			},
		},
		"post path": {
			newRequest("POST", "http://example.com/path/here"),
			map[string]interface{}{
				"http.method":     "POST",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/path/here",
				"http.path":       "/path/here",
				"http.host":       "example.com",
				"http.user":       "",
				"http.ref":        "",
				"http.user-agent": "",
			},
			map[string]string{
				"transaction": `(?:[0-9a-z]+-){4}[0-9a-z]+`,
			},
		},
		"strips params off method": {
			newRequest("GET", "http://example.com/token/1/test?apid=1&thing=2"),
			map[string]interface{}{
				"http.method":     "GET",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/token/1/test?apid=1&thing=2",
				"http.path":       "/token/1/test",
				"http.host":       "example.com",
				"http.user":       "",
				"http.ref":        "",
				"http.user-agent": "",
			},
			map[string]string{
				"transaction": `(?:[0-9a-z]+-){4}[0-9a-z]+`,
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
				RemoteAddr: "192.168.100.5",
			},
			map[string]interface{}{
				"http.method":     "CONNECT",
				"http.protocol":   "HTTP/2.0",
				"http.uri":        "www.example.com:443",
				"http.path":       "www.example.com:443",
				"http.host":       "www.example.com:443",
				"http.user":       "",
				"http.ref":        "",
				"http.user-agent": "",
			},
			map[string]string{
				"transaction": `(?:[0-9a-z]+-){4}[0-9a-z]+`,
			},
		},
		"user": {
			userRequest,
			map[string]interface{}{
				"http.method":     "GET",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/",
				"http.path":       "/",
				"http.host":       "example.com",
				"http.user":       "192.168.100.5",
				"http.ref":        "",
				"http.user-agent": "",
			},
			map[string]string{
				"transaction": `(?:[0-9a-z]+-){4}[0-9a-z]+`,
			},
		},
		"referer": {
			referrerRequest,
			map[string]interface{}{
				"http.method":     "GET",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/test",
				"http.path":       "/test",
				"http.host":       "example.com",
				"http.user":       "",
				"http.ref":        "http://google.com",
				"http.user-agent": "",
			},
			map[string]string{
				"transaction": `(?:[0-9a-z]+-){4}[0-9a-z]+`,
			},
		},
		"user agent": {
			userAgentRequest,
			map[string]interface{}{
				"http.method":     "GET",
				"http.protocol":   "HTTP/1.1",
				"http.uri":        "/",
				"http.path":       "/",
				"http.host":       "example.com",
				"http.user":       "",
				"http.ref":        "",
				"http.user-agent": "some user agent",
			},
			map[string]string{
				"transaction": `(?:[0-9a-z]+-){4}[0-9a-z]+`,
			},
		},
	}

	rec := httptest.NewRecorder()

	for k, tc := range cases {
		beforeHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			entry := log.Ctx(req.Context())
			for f, v := range tc.expected {
				assert.Contains(t, entry.Data, f, "test %s - Has Field: %s", k, f)
				assert.Equal(t, v, entry.Data[f], "test %s - Field: %s", k, f)
			}
			for f, v := range tc.regex {
				assert.Contains(t, entry.Data, f, "test %s - Has Field: %s", k, f)
				assert.Regexp(t, v, entry.Data[f], "test %s - Field: %s", k, f)
			}
			w.Write([]byte("ok\n"))
		})

		handler := LogContextHandler(beforeHandler)
		handler.ServeHTTP(rec, tc.request)
	}
}
