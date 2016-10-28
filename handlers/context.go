// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package handlers

import (
	"net/http"

	log "github.com/graze/golang-service/log"
)

type logContextHandler struct {
	logger  log.LogContext
	handler http.Handler
}

// ServeHTTP does the actual handling of HTTP requests by wrapping the request in a logger
func (h logContextHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := MakeLogger(w, h.logger)
	url := *req.URL
	context := logger.GetContext()
	context.Add(log.F{
		"http.method":   req.Method,
		"http.protocol": req.Proto,
		"http.uri":      parseUri(req, url),
		"http.path":     uriPath(req, url),
		"http.host":     req.Host,
	})
	h.handler.ServeHTTP(logger, req)
}

// LogContextHandler returns a handler that adds `http` and `transaction` items into a common logging context
//
// It adds the following fields to the `LoggingResponseWriter` log context:
//  http.method     - GET/POST/...
//  http.protocol   - HTTP/1.1
//  http.uri        - /path?with=query
//  http.path       - /path
//  http.host       - localhost:80
//  transaction     - unique uuid4 for this request
func LogContextHandler(h http.Handler) http.Handler {
	return logContextHandler{log.New(), h}
}

// LogContextHandler returns a handler that adds `http` and `transaction` items into the provided logging context
func LoggingContextHandler(logger log.LogContext, h http.Handler) http.Handler {
	return logContextHandler{logger.With(log.F{}), h}
}
