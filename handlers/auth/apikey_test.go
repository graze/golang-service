// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok\n"))
})

// headerRequest creates a new request with the authHeader
func headerRequest(t *testing.T, method, url string, headers map[string]string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return req
}

func TestApiKeyAuthErrors(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		provider string
		request  *http.Request
		err      error
		status   int
		finder   UserFinder
	}{
		"no header": {
			"Graze",
			headerRequest(t, "GET", "/path", map[string]string{}),
			&NoHeaderError{},
			http.StatusUnauthorized,
			func(key string, r *http.Request) (interface{}, error) {
				return "", nil
			},
		},
		"invalid provider": {
			"Graze",
			headerRequest(t, "GET", "/path", map[string]string{"Authorization": "Fish cake"}),
			&BadProviderError{"Graze", "Fish"},
			http.StatusUnauthorized,
			func(key string, r *http.Request) (interface{}, error) {
				return "", nil
			},
		},
		"invalid format": {
			"Graze",
			headerRequest(t, "GET", "/path", map[string]string{"Authorization": "Fish"}),
			&InvalidFormatError{"<provider> <apiKey>", "Fish"},
			http.StatusUnauthorized,
			func(key string, r *http.Request) (interface{}, error) {
				return "", nil
			},
		},
		"invalid format - too many fields": {
			"Graze",
			headerRequest(t, "GET", "/path", map[string]string{"Authorization": "Fish cake thing"}),
			&InvalidFormatError{"<provider> <apiKey>", "Fish cake thing"},
			http.StatusUnauthorized,
			func(key string, r *http.Request) (interface{}, error) {
				return "", nil
			},
		},
		"failed finder": {
			"Graze",
			headerRequest(t, "GET", "/path", map[string]string{"Authorization": "Graze key"}),
			&InvalidKeyError{"key", errors.New("")},
			http.StatusUnauthorized,
			func(key string, r *http.Request) (interface{}, error) {
				assert.Equal(t, "key", key)
				return "", errors.New("some failed error")
			},
		},
	}

	rec := httptest.NewRecorder()

	for k, tc := range cases {
		auth := &APIKey{tc.provider, tc.finder, func(w http.ResponseWriter, r *http.Request, err error, status int) {
			assert.IsType(t, tc.err, err, "test: %s", k)
			assert.Equal(t, tc.status, status, "test: %s", k)
		}}
		handler := auth.Then(okHandler)
		handler.ServeHTTP(rec, tc.request)
	}
}

func TestUserStorage(t *testing.T) {
	t.Parallel()

	user := "some user"

	cases := map[string]struct {
		request  *http.Request
		provider string
		finder   UserFinder
		expected interface{}
	}{
		"nil return": {
			headerRequest(t, "GET", "/stuff", map[string]string{"Authorization": "Graze key"}),
			"Graze",
			func(key string, r *http.Request) (interface{}, error) {
				assert.Equal(t, "key", key)
				return nil, nil
			},
			interface{}(nil),
		},
		"user": {
			headerRequest(t, "GET", "/stuff", map[string]string{"Authorization": "Graze otherKey"}),
			"Graze",
			func(key string, r *http.Request) (interface{}, error) {
				assert.Equal(t, "otherKey", key)
				return user, nil
			},
			user,
		},
	}

	rec := httptest.NewRecorder()

	for k, tc := range cases {
		auth := &APIKey{tc.provider, tc.finder, func(w http.ResponseWriter, r *http.Request, err error, status int) {
			t.Errorf("onError handler called. Err: %s, Status: %d, Test: %s", err, status, k)
		}}

		baseHandler := func(w http.ResponseWriter, req *http.Request) {
			user := GetUser(req)
			assert.Equal(t, tc.expected, user, "test: %s", k)
		}

		handler := auth.ThenFunc(baseHandler)
		handler.ServeHTTP(rec, tc.request)
	}
}
