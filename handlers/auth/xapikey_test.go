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

	"github.com/graze/golang-service/handlers/failure"
	"github.com/stretchr/testify/assert"
)

func TestXApiKeyAuthErrors(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		request *http.Request
		err     error
		status  int
		finder  Finder
	}{
		"no header": {
			headerRequest(t, "GET", "/path", map[string]string{}),
			&NoHeaderError{},
			http.StatusUnauthorized,
			FinderFunc(func(key interface{}, r *http.Request) (interface{}, error) {
				return "", nil
			}),
		},
		"failed finder": {
			headerRequest(t, "GET", "/path", map[string]string{"x-api-key": "key"}),
			&InvalidKeyError{"key", errors.New("")},
			http.StatusUnauthorized,
			FinderFunc(func(key interface{}, r *http.Request) (interface{}, error) {
				assert.Equal(t, "key", key)
				return "", errors.New("some failed error")
			}),
		},
	}

	rec := httptest.NewRecorder()

	for k, tc := range cases {
		auth := NewXAPIKey(tc.finder, failure.HandlerFunc(func(w http.ResponseWriter, r *http.Request, err error, status int) {
			assert.IsType(t, tc.err, err, "test: %s", k)
			assert.Equal(t, tc.status, status, "test: %s", k)
		}))
		handler := auth.Then(okHandler)
		handler.ServeHTTP(rec, tc.request)
	}
}

func TestValidXAPIKeyAuth(t *testing.T) {
	t.Parallel()

	user := "some user"

	cases := map[string]struct {
		request  *http.Request
		finder   Finder
		expected interface{}
	}{
		"nil return": {
			headerRequest(t, "GET", "/stuff", map[string]string{"X-API-Key": "key"}),
			FinderFunc(func(key interface{}, r *http.Request) (interface{}, error) {
				assert.Equal(t, "key", key)
				return nil, nil
			}),
			interface{}(nil),
		},
		"user": {
			headerRequest(t, "GET", "/stuff", map[string]string{"X-api-Key": "otherKey"}),
			FinderFunc(func(key interface{}, r *http.Request) (interface{}, error) {
				assert.Equal(t, "otherKey", key)
				return user, nil
			}),
			user,
		},
	}

	rec := httptest.NewRecorder()

	for k, tc := range cases {
		auth := NewXAPIKey(tc.finder, failure.HandlerFunc(func(w http.ResponseWriter, r *http.Request, err error, status int) {
			t.Errorf("onError handler called. Err: %s, Status: %d, Test: %s", err, status, k)
		}))

		baseHandler := func(w http.ResponseWriter, req *http.Request) {
			user := GetUser(req)
			assert.Equal(t, tc.expected, user, "test: %s", k)
		}

		handler := auth.ThenFunc(baseHandler)
		handler.ServeHTTP(rec, tc.request)
	}
}
