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
	"errors"
	"io"
	"net/http"
)

// Handler handlers a panic recovery and does something (outputs to w, logs, reports to third party, etc)
//
// Note that multiple Recoverers could write to w
type Handler interface {
	Handle(w io.Writer, r *http.Request, err error, status int)
}

// HandlerFunc provides a simple function to handle when a http.Handler panic occours
type HandlerFunc func(io.Writer, *http.Request, error, int)

// Handle implements the Handler interface for a HandlerFunc
func (f HandlerFunc) Handle(w io.Writer, r *http.Request, err error, status int) {
	f(w, r, err, status)
}

// middleware is a http.Handler to recover from panics
type middleware struct {
	next     http.Handler
	handlers []Handler
}

// ServeHTTP defers a panic handler and writes 500, then passes off the error to something else
func (h middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			err, ok := e.(error)
			if !ok {
				err = errors.New(e.(string))
			}

			w.WriteHeader(http.StatusInternalServerError)
			for _, r := range h.handlers {
				r.Handle(w, req, err, http.StatusInternalServerError)
			}
		}
	}()

	h.next.ServeHTTP(w, req)
}

// New creates a http.Handler middleware that loops through a series
//
// Usage:
// 	r := mux.NewRouter()
// 	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 	   panic("oh-o")
// 	})
//
// 	outputRecoverer := func(w io.Writer, r *http.Request, err error, status int) {
// 		w.Write([]byte("panic happened, oh dear"))
// 	}
// 	recoverer := recovery.New(r, recovery.Logger(log.New()), raygun.New(raygunClient), recovery.HandlerFunc(format))
// 	http.ListenAndServe(":80", recoverer)
func New(h http.Handler, handlers ...Handler) http.Handler {
	return middleware{h, handlers}
}
