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

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestContextUpdatesTheRequestContext(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		request *http.Request
		before  map[string]interface{}
		after   map[string]string
	}{
		"basic": {
			newRequest("GET", "http://example.com"),
			map[string]interface{}{
				"http.method":   "GET",
				"http.protocol": "HTTP/1.1",
				"http.uri":      "/",
				"http.path":     "/",
				"http.host":     "example.com",
			},
			map[string]string{
				"http.status": "200",
				"http.bytes":  `\d+`,
				"http.dur":    `[0-9\.]+`,
				"transaction": `(?:[0-9a-z]+-){4}[0-9a-z]+`,
			},
		},
		"post path": {
			newRequest("POST", "http://example.com/path/here"),
			map[string]interface{}{
				"http.method":   "POST",
				"http.protocol": "HTTP/1.1",
				"http.uri":      "/path/here",
				"http.path":     "/path/here",
				"http.host":     "example.com",
			},
			map[string]string{
				"http.status": "200",
				"http.bytes":  `\d+`,
				"http.dur":    `[0-9\.]+`,
				"transaction": `(?:[0-9a-z]+-){4}[0-9a-z]+`,
			},
		},
		"strips params off method": {
			newRequest("GET", "http://example.com/token/1/test?apid=1&thing=2"),
			map[string]interface{}{
				"http.method":   "GET",
				"http.protocol": "HTTP/1.1",
				"http.uri":      "/token/1/test?apid=1&thing=2",
				"http.path":     "/token/1/test",
				"http.host":     "example.com",
			},
			map[string]string{
				"http.status": "200",
				"http.bytes":  `\d+`,
				"http.dur":    `[0-9\.]+`,
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
				"http.method":   "CONNECT",
				"http.protocol": "HTTP/2.0",
				"http.uri":      "www.example.com:443",
				"http.path":     "www.example.com:443",
				"http.host":     "www.example.com:443",
			},
			map[string]string{
				"http.status": "200",
				"http.bytes":  `\d+`,
				"http.dur":    `[0-9\.]+`,
				"transaction": `(?:[0-9a-z]+-){4}[0-9a-z]+`,
			},
		},
	}

	for k, tc := range cases {
		var logger LoggingResponseWriter = nil
		beforeHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if _, ok := w.(LoggingResponseWriter); ok {
				logger = w.(LoggingResponseWriter)
				if _, ok := logger.GetContext().(*logrus.Entry); ok {
					entry := logger.GetContext().(*logrus.Entry)
					for f, v := range tc.before {
						assert.Contains(t, entry.Data, f, "test %s - Has Field: %s", k, f)
						assert.Equal(t, v, entry.Data[f], "test %s - Field: %s", k, f)
					}
				} else {
					t.Error("returned context does not implement logrus.Entry so unable to retrieve data")
				}
			} else {
				t.Error("http.ResponseWriter should implement LoggingResponseWriter")
			}
			w.Write([]byte("ok\n"))
		})

		handler := LogContextHandler(beforeHandler)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, tc.request)
		if _, ok := logger.GetContext().(*logrus.Entry); ok {
			entry := logger.GetContext().(*logrus.Entry)
			for f, v := range tc.after {
				assert.Contains(t, entry.Data, f, "test %s - Has Field: %s", k, f)
				assert.Regexp(t, v, entry.Data[f], "test %s - Field: %s", k, f)
			}
		} else {
			t.Error("returned context does not implement logrus.Entry so unable to retrieve data")
		}
	}
}