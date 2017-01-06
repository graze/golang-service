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
	"fmt"
	"net/http"
	"strings"
)

// FailHandler gets called if the handler found an error with the Authorization
type FailHandler func(w http.ResponseWriter, r *http.Request, err error, status int)

// APIKey contains a wrapper around a handler to provide authentication
//
// It uses the Authorization header in the format: <provider> <apiKey>
// If the format of the header is valid, the validator will be called with the apiKey
// if anything goes wrong, a callback on onError is called with the error and the http StatusCode to return
type APIKey struct {
	// Provider is the name of the key being provided. The Authorization header is in the format: <provider> <apiKey>
	// It must not contain any spaces
	Provider string
	// Validator takes the provided <apiKey> and returns a user object or error if the key is invalid
	Finder Finder
	// OnError gets called if the request is unauthorized or forbidden
	OnError FailHandler
}

// apiKeyHandler implements http.Handler and responds
type apiKeyHandler struct {
	apiKey  APIKey
	handler http.Handler
}

type (
	// NoHeaderError for when the Authorization header is not provided
	NoHeaderError struct{}
	// InvalidFormatError if the Authorization header is not in the format: <provider> <apiKey>
	InvalidFormatError struct{ format, header string }
	// BadProviderError when the supplied provider does not match the expected
	BadProviderError struct{ provider, expected string }
	// InvalidKeyError if the supplied key does not match any existing keys
	InvalidKeyError struct {
		key string
		err error
	}
)

func (e *NoHeaderError) Error() string {
	return "no Authorization header provided"
}

func (e *InvalidFormatError) Error() string {
	return fmt.Sprintf("provided Authorization header in invalid format, expecting: %s got: %s", e.format, e.header)
}

func (e *BadProviderError) Error() string {
	return fmt.Sprintf("Authroziation provider does not match. Expecting: %s got: %s", e.expected, e.provider)
}

func (e *InvalidKeyError) Error() string {
	return fmt.Sprintf("provided api key: '%s' is not valid: %s", e.key, e.err.Error())
}

// ThenFunc surrounds an existing handler func and returns a new http.Handler
//
// Usage:
//  func finder(creds interface{}, r *http.Request) (interface{}, error) {
// 		key, ok := creds.(string)
// 		if !ok {
// 			return nil, fmt.Errorf("Could not understand creds")
// 		}
//      user, ok := users[key]
//      if !ok {
//          return nil, fmt.Errorf("No user found for: %s", key)
//      }
//      return user, nil
//  }
//
//  func onError(w http.ResponseWriter, r *http.Request, err error, status int) {
//      w.WriteHeader(status)
//      fmt.Fprintf(w, err.Error())
//  }
//
//  keyAuth := auth.APIKey{"Graze", finder, onError}
//
//  http.Handle("/thing", keyAuth.ThenFunc(ThingFunc))
func (w APIKey) ThenFunc(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return w.Then(http.HandlerFunc(fn))
}

// Then surrounds an existing http.Handler and returns a new http.Handler
//
// Usage:
//  func finder(creds interface{}, r *http.Request) (interface{}, error) {
// 		key, ok := creds.(string)
// 		if !ok {
// 			return nil, fmt.Errorf("Could not understand creds")
// 		}
//      user, ok := users[key]
//      if !ok {
//          return nil, fmt.Errorf("No user found for: %s", key)
//      }
//      return user, nil
//  }
//
//  func onError(w http.ResponseWriter, r *http.Request, err error, status int) {
//      w.WriteHeader(status)
//      fmt.Fprintf(w, err.Error())
//  }
//
//  keyAuth := auth.APIKey{"Graze", finder, onError}
//
//  http.Handle("/thing", keyAuth.Then(ThingHandler))
func (w APIKey) Then(h http.Handler) http.Handler {
	return apiKeyHandler{w, h}
}

// Handler wraps the Then method to become clearer
func (w APIKey) Handler(h http.Handler) http.Handler {
	return w.Then(h)
}

// ServeHTTP checks if the request has the correct authentication
func (h apiKeyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	header := req.Header["Authorization"]
	if len(header) == 0 {
		h.apiKey.OnError(w, req, &NoHeaderError{}, http.StatusUnauthorized)
		return
	}

	parts := strings.Split(header[0], " ")
	if len(parts) != 2 {
		h.apiKey.OnError(w, req, &InvalidFormatError{"<provider> <apiKey>", header[0]}, http.StatusUnauthorized)
		return
	}

	provider, value := parts[0], parts[1]
	if provider != h.apiKey.Provider {
		h.apiKey.OnError(w, req, &BadProviderError{provider, h.apiKey.Provider}, http.StatusUnauthorized)
		return
	}

	user, err := h.apiKey.Finder.Find(value, req)
	if err != nil {
		h.apiKey.OnError(w, req, &InvalidKeyError{value, err}, http.StatusUnauthorized)
		return
	}
	req = saveUser(req, user)

	h.handler.ServeHTTP(w, req)
}
