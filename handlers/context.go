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
	"time"

	log "github.com/graze/golang-service/logging"
)

type logContextHandler struct {
	logger  log.LogContext
	handler http.Handler
}

// ServeHTTP does the actual handling of HTTP requests by wrapping the request in a logger
func (h logContextHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := time.Now().UTC()
	logger := MakeLogger(w, h.logger)
	url := *req.URL
	context := logger.GetContext()
	context.Info("test of MakeLogger GetContext()")
	h.logger.Info("test of logContextHandler log context")
	context.Add(log.F{
		"http.method":   req.Method,
		"http.protocol": req.Proto,
		"http.uri":      parseUri(req, url),
		"http.path":     uriPath(req, url),
		"http.host":     req.Host,
	})
	h.handler.ServeHTTP(logger, req)
	dur := time.Now().UTC().Sub(t)
	sDur := float64(dur.Nanoseconds()) / (float64(time.Second) / float64(time.Nanosecond))
	context.Add(log.F{
		"http.status": logger.Status(),
		"http.bytes":  logger.Size(),
		"http.dur":    sDur,
	})
}

// LogContextHandler returns a handler that adds `http` and `transaction` items into a common logging context
func LogContextHandler(h http.Handler) http.Handler {
	return logContextHandler{log.New(), h}
}

// LogContextHandler returns a handler that adds `http` and `transaction` items into the provided logging context
func LoggingContextHandler(logger log.LogContext, h http.Handler) http.Handler {
	return logContextHandler{logger.With(log.F{}), h}
}
