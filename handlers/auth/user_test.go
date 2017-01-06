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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserStorage(t *testing.T) {
	t.Parallel()

	user := "some user"

	cases := map[string]struct {
		request  *http.Request
		provider string
		finder   Finder
		expected interface{}
	}{
		"nil return": {
			headerRequest(t, "GET", "/stuff", map[string]string{"Authorization": "Graze key"}),
			"Graze",
			FinderFunc(func(key interface{}, r *http.Request) (interface{}, error) {
				assert.Equal(t, "key", key)
				return nil, nil
			}),
			interface{}(nil),
		},
		"user": {
			headerRequest(t, "GET", "/stuff", map[string]string{"Authorization": "Graze otherKey"}),
			"Graze",
			FinderFunc(func(key interface{}, r *http.Request) (interface{}, error) {
				assert.Equal(t, "otherKey", key)
				return user, nil
			}),
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
