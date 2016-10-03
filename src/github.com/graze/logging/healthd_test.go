// This file is part of graze/golang-service-logging
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// @license https://github.com/graze/golang-service-logging/blob/master/LICENSE
// @link    https://github.com/graze/golang-service-logging
package logging

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "time"
    "net/http"
    "bytes"
    "strings"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok\n"))
})

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestHealthdLogging(t *testing.T) {
    loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Fatal(err)
	}
	ts := time.Date(1983, 05, 26, 3, 30, 45, int((736 * time.Millisecond).Nanoseconds()), loc)

	// A typical request with an OK response
	req := newRequest("GET", "http://example.com")
    req.Header.Add("X-Forwarded-For", "192.168.100.5")

    buf := new(bytes.Buffer)
    dur, err := time.ParseDuration("0.302s")
    if (err != nil) {
        t.Fatal(err)
    }
    writeHealthdLog(buf, req, *req.URL, ts, dur, http.StatusOK, 100)
    log := buf.String()

    assert.Equal(t, strings.Join([]string{`422764245.736"/"200"0.302"0.302"192.168.100.5`,"\n"}, ""), log)

    ts = time.Date(1983, 05, 26, 3, 30, 45, int((123 * time.Millisecond).Nanoseconds()), loc)
    req = newRequest("GET", "http://example.com/path/here")

    buf = new(bytes.Buffer)
    dur, err = time.ParseDuration("0.102s")
    if (err != nil) {
        t.Fatal(err)
    }
    writeHealthdLog(buf, req, *req.URL, ts, dur, http.StatusOK, 100)
    log = buf.String()

    assert.Equal(t, strings.Join([]string{`422764245.123"/path/here"200"0.102"0.102"`,"\n"}, ""), log)
}
