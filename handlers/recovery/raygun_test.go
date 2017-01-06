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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MindscapeHQ/raygun4go"
	"github.com/graze/golang-service/log"
	"github.com/stretchr/testify/assert"
)

type raygunMock struct {
	request *http.Request
	data    interface{}
	err     string
}

func (r *raygunMock) Request(req *http.Request) *raygun4go.Client {
	r.request = req
	return nil
}

func (r *raygunMock) CustomData(data interface{}) *raygun4go.Client {
	r.data = data
	return nil
}

func (r *raygunMock) CreateError(message string) error {
	r.err = message
	return nil
}

func TestRaygun(t *testing.T) {
	cases := map[string]struct {
		req  *http.Request
		data interface{}
		err  string
	}{
		"base": {
			newRequest("GET", "http://example.com").WithContext(log.With(log.KV{"key": "value"}).NewContext(context.Background())),
			log.KV{"key": "value"},
			"oh no!",
		},
	}

	for k, tc := range cases {
		rec := httptest.NewRecorder()
		mock := &raygunMock{}
		handler := New(Raygun(mock)).Handle(panicHandler)
		handler.ServeHTTP(rec, tc.req)

		assert.Equal(t, tc.req, mock.request, "test: %s", k)
		assert.Equal(t, tc.data, mock.data, "test: %s", k)
		assert.Equal(t, tc.err, mock.err, "test: %s", k)
	}
}
