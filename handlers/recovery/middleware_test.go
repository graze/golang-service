// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package recovery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/graze/golang-service/handlers/failure"
	"github.com/stretchr/testify/assert"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok\n"))
})

var panicHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	panic("oh no!")
})

var echoRecoverer = failure.HandlerFunc(func(w http.ResponseWriter, r *http.Request, err error, status int) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
})

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestHandlerCallsNextHandlerWhenNoPanicOccours(t *testing.T) {
	handler := New()(okHandler)

	rec := httptest.NewRecorder()
	req := newRequest("GET", "http://example.com")

	handler.ServeHTTP(rec, req)

	assert.Equal(t, "ok\n", rec.Body.String())
}

func TestPanics(t *testing.T) {
	cases := map[string]struct {
		handlers []failure.Handler
		body     string
		status   int
	}{
		"echo": {
			[]failure.Handler{echoRecoverer},
			"oh no!",
			http.StatusInternalServerError,
		},
		"multiple": {
			[]failure.Handler{echoRecoverer, echoRecoverer},
			"oh no!oh no!",
			http.StatusInternalServerError,
		},
	}

	for k, tc := range cases {
		rec := httptest.NewRecorder()
		handler := New(tc.handlers...)(panicHandler)
		handler.ServeHTTP(rec, newRequest("GET", "http://example.com"))
		assert.Equal(t, tc.body, rec.Body.String(), "test: %s", k)
		assert.Equal(t, tc.status, rec.Code, "test: %s", k)
	}
}
