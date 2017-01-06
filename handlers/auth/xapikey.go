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

import "net/http"

// XAPIKey contains a wrapper around a handler to provide authentication using the X-Api-Key header
//
// It uses the x-api-key header in the format: <apiKey>
// if anything goes wrong, a callback on onError is called with the error and the http StatusCode to return
type XAPIKey struct {
	// Validator takes the provided <apiKey> and returns a user object or error if the key is invalid
	Finder Finder
	// OnError gets called if the request is unauthorized or forbidden
	OnError FailHandler
}

// ThenFunc surrounds an existing handler func and returns a new http.Handler
//
// Usage:
//  func finder(creds interface{}, r *http.Request) (interface{}, error) {
// 		key, ok := creds.(string)
// 		if !ok {
// 			return nil, fmt.Errorf("Could not understand creds")
// 		}
// 		user, ok := users[key]
// 		if !ok {
// 			return nil, fmt.Errorf("No user found for: %s", key)
// 		}
// 		return user, nil
// 	}
//
// 	func onError(w http.ResponseWriter, r *http.Request, err error, status int) {
// 		w.WriteHeader(status)
// 		fmt.Fprintf(w, err.Error())
// 	}
//
// 	keyAuth := auth.APIKey{"Graze", finder, onError}
//
// 	http.Handle("/thing", keyAuth.ThenFunc(ThingFunc))
func (x XAPIKey) ThenFunc(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return x.Handler(http.HandlerFunc(fn))
}

// Then surrounds an existing http.Handler and returns a new http.Handler
//
// Usage:
// 	func finder(creds interface{}, r *http.Request) (interface{}, error) {
// 		key, ok := creds.(string)
// 		if !ok {
// 			return nil, fmt.Errorf("Could not understand creds")
// 		}
// 		user, ok := users[key]
// 		if !ok {
// 			return nil, fmt.Errorf("No user found for: %s", key)
// 		}
// 		return user, nil
// 	}
//
// 	func onError(w http.ResponseWriter, r *http.Request, err error, status int) {
// 		w.WriteHeader(status)
// 		fmt.Fprintf(w, err.Error())
// 	}
//
// 	keyAuth := auth.APIKey{"Graze", finder, onError}
//
// 	http.Handle("/thing", keyAuth.Then(ThingHandler))
func (x XAPIKey) Then(h http.Handler) http.Handler {
	return x.Handler(h)
}

// Handler wraps the Then method to become clearer
func (x XAPIKey) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		header := req.Header["X-Api-Key"]
		if len(header) == 0 {
			x.OnError(w, req, &NoHeaderError{}, http.StatusUnauthorized)
			return
		}

		user, err := x.Finder.Find(header[0], req)
		if err != nil {
			x.OnError(w, req, &InvalidKeyError{header[0], err}, http.StatusUnauthorized)
			return
		}
		req = saveUser(req, user)

		h.ServeHTTP(w, req)
	})
}

// NewXAPIKey returns an APIKey struct that has a Handle method to provide authentication to your service
func NewXAPIKey(finder Finder, onError FailHandler) *XAPIKey {
	return &XAPIKey{finder, onError}
}
