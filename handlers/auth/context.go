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
	"context"
	"net/http"
)

// contextKey is a custom type to only allow this to access the key in the context
type contextKey int

// userKey is a private key to store the user information in the context in
const userKey contextKey = iota

// SaveUser takes a nominal user and stores it in a new context for the provided request
func SaveUser(r *http.Request, user interface{}) *http.Request {
	if user == nil {
		return r
	}
	return r.WithContext(context.WithValue(r.Context(), userKey, user))
}

// GetUser retrieves any user information provided by the validate request
//
// This is to handle different apiKeys having potentially different permissions
// or having different backend resources
//
// Usage:
// 	keyAuth := auth.ApiKey{"Graze", finder, onError}
//
// 	http.Handle("/thing", keyAuth.Handler(ItemHandler))
//
// 	func ItemHandler(w http.ResponseWriter, r *http.Request) {
// 		user, ok := auth.GetUser(r.Context()).(*User)
// 		if !ok {
// 			w.WriteHeader(403)
// 			return
// 		}
// 		...
// 	}
func GetUser(r *http.Request) interface{} {
	return r.Context().Value(userKey)
}
