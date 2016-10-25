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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestHealthdFileLogging(t *testing.T) {
	handler := HealthdHandler(okHandler)

	// A typical request with an OK response
	req := newRequest("GET", "http://example.com/")

	rec := httptest.NewRecorder()
	timestamp := time.Now().UTC().Format("2006-01-02-15")
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	file := "/var/log/nginx/healthd/application.log." + timestamp

	fe, err := exists("/var/log/nginx/healthd")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, fe)
	fe, err = exists(file)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, fe)

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, `[0-9\.]+"/"200"[0-9\.]+"[0-9\.]+"[0-9\.]*`, string(bytes))
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func TestHealthdLogging(t *testing.T) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatal(err)
	}

	headerRequest := newRequest("GET", "http://example.com")
	headerRequest.Header.Add("X-Forwarded-For", "192.168.100.5")

	cases := map[string]struct {
		ts       time.Time
		dur      time.Duration
		req      *http.Request
		expected string
	}{
		"with header": {
			time.Date(1983, 05, 26, 3, 30, 45, int((736 * time.Millisecond).Nanoseconds()), loc),
			getDuration(t, "0.302s"),
			headerRequest,
			`422767845.736"/"200"0.302"0.302"192.168.100.5`,
		},
		"standard": {
			time.Date(1983, 05, 26, 3, 30, 45, int((123 * time.Millisecond).Nanoseconds()), loc),
			getDuration(t, "0.102s"),
			newRequest("POST", "http://example.com/path/here"),
			`422767845.123"/path/here"200"0.102"0.102"`,
		},
	}

	for k, tc := range cases {
		buf := new(bytes.Buffer)
		writeHealthdLog(buf, tc.req, *tc.req.URL, tc.ts, tc.dur, http.StatusOK, 100)
		log := buf.String()

		assert.Equal(t, tc.expected+"\n", log, "test: %s", k)
	}
}
