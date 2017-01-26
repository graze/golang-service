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
	"net/http"

	"github.com/graze/golang-service/handlers/failure"
)

// middleware provides a Handle method that implements http.Handler, and can be used with other http handler middlewares
type middleware struct {
	next     http.Handler
	handlers []failure.Handler
}

// Handle returns a middleware http.Handler to be used when handling requests
func (m *middleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			err, ok := e.(error)
			if !ok {
				err = errors.New(e.(string))
			}

			for _, r := range m.handlers {
				r.Handle(w, req, err, http.StatusInternalServerError)
			}
		}
	}()

	m.next.ServeHTTP(w, req)
}

// New creates a http.Handler middleware that loops through a series of panic handlers that can write data back to the
// response or log the panic
//
// The handler will always write a header of 500 (Internal Server Error) and each Panic handler can add content to the
// body if required
//
// Usage:
// 	r := mux.NewRouter()
// 	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 	   panic("oh-o")
// 	})
//
// 	outputRecoverer := func(w io.Writer, r *http.Request, err error, status int) {
// 		w.Write([]byte(`{"error":"unknown error"}`))
// 	}
// 	recoverer := recovery.New(panic.Logger(log.New()), raygun.New(raygunClient), panic.HandlerFunc(format))
// 	http.ListenAndServe(":80", recoverer(r))
func New(handlers ...failure.Handler) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &middleware{h, handlers}
	}
}
