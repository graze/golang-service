// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

/*
Package failure provides a generic handler to do something when it all goes wrong during an http request

Sometimes bad things happen and you need to handle them nicely

This interface is used within this library in the auth and recovery packages.

    type Handler interface {
        Handle(w http.ResponseWriter, r *http.Request, err error, status int)
    }

The HandlerFunc converts a function to a Handler interface

    errorHandler = failure.Handler(func(w http.ResponseWriter, r *http.Request, err error, status int) {
        w.WriteHeader(status)
        w.Write([]byte(err.Error()))
    })
*/
package failure

import "net/http"

// Handler handlers a panic panic and does something (outputs to w, logs, reports to third party, etc)
//
// Note that multiple Recoverers could write to w
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request, err error, status int)
}

// HandlerFunc provides a simple function to handle when a http.Handler panic occours
type HandlerFunc func(http.ResponseWriter, *http.Request, error, int)

// Handle implements the Handler interface for a HandlerFunc
func (f HandlerFunc) Handle(w http.ResponseWriter, r *http.Request, err error, status int) {
	f(w, r, err, status)
}
