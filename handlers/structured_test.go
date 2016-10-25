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
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestStructuredLogging(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()

	headerRequest := newRequest("GET", "http://example.com")
	headerRequest.Header.Add("X-Forwarded-For", "192.168.100.5")

	referrerRequest := newRequest("GET", "http://example.com/test")
	referrerRequest.Header.Add("Referer", "http://google.com")

	cases := map[string]struct {
		request   *http.Request
		timestamp time.Time
		duration  time.Duration
		size      int
		expected  string
	}{
		"basic": {
			newRequest("GET", "http://example.com"),
			now,
			getDuration(t, "0.302s"),
			100,
			fmt.Sprintf(
				`tag=request_handled msg="GET / HTTP/1.1" method=GET protocol=HTTP/1.1 uri=/ path=/ host=example.com status=200 bytes=100 dur=0.302 ts=%s ref= user=`+"\n",
				now.Format(time.RFC3339Nano)),
		},
		"post path": {
			newRequest("POST", "http://example.com/path/here"),
			now,
			getDuration(t, "0.102s"),
			200,
			fmt.Sprintf(
				`tag=request_handled msg="POST /path/here HTTP/1.1" method=POST protocol=HTTP/1.1 uri=/path/here path=/path/here host=example.com status=200 bytes=200 dur=0.102 ts=%s ref= user=`+"\n",
				now.Format(time.RFC3339Nano)),
		},
		"strips params off method": {
			newRequest("GET", "http://example.com/token/1/test?apid=1&thing=2"),
			now,
			getDuration(t, "0.927s"),
			300,
			fmt.Sprintf(
				`tag=request_handled msg="GET /token/1/test?apid=1&thing=2 HTTP/1.1" method=GET protocol=HTTP/1.1 uri="/token/1/test?apid=1&thing=2" path=/token/1/test host=example.com status=200 bytes=300 dur=0.927 ts=%s ref= user=`+"\n",
				now.Format(time.RFC3339Nano)),
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
			now,
			getDuration(t, "0.927s"),
			400,
			fmt.Sprintf(
				`tag=request_handled msg="CONNECT www.example.com:443 HTTP/2.0" method=CONNECT protocol=HTTP/2.0 uri=www.example.com:443 path=www.example.com:443 host=www.example.com:443 status=200 bytes=400 dur=0.927 ts=%s ref= user=`+"\n",
				now.Format(time.RFC3339Nano)),
		},
		"handles referrer": {
			referrerRequest,
			now,
			getDuration(t, "0.019s"),
			500,
			fmt.Sprintf(
				`tag=request_handled msg="GET /test HTTP/1.1" method=GET protocol=HTTP/1.1 uri=/test path=/test host=example.com status=200 bytes=500 dur=0.019 ts=%s ref=http://google.com user=`+"\n",
				now.Format(time.RFC3339Nano)),
		},
		"handles x-forwarded-for": {
			headerRequest,
			now,
			getDuration(t, "0.019s"),
			600,
			fmt.Sprintf(
				`tag=request_handled msg="GET / HTTP/1.1" method=GET protocol=HTTP/1.1 uri=/ path=/ host=example.com status=200 bytes=600 dur=0.019 ts=%s ref= user=192.168.100.5`+"\n",
				now.Format(time.RFC3339Nano)),
		},
	}

	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	for k, tc := range cases {
		buf.Reset()
		writeStructuredLog(logger, tc.request, *tc.request.URL, tc.timestamp, tc.duration, http.StatusOK, tc.size)
		assert.Equal(t, tc.expected, buf.String(), "test: %s", k)
	}
}
