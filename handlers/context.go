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

	"github.com/Sirupsen/logrus"
)

type logContextHandler struct {
	handler http.Handler
}

// ServeHTTP does the actual handling of HTTP requests by wrapping the request in a logger
func (h logContextHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := time.Now().UTC()
	logger := MakeLogger(w)
	url := *req.URL
	logger.AddContext(logrus.Fields{
		"http.method":   req.Method,
		"http.protocol": req.Proto,
		"http.uri":      parseUri(req, url),
		"http.path":     uriPath(req, url),
		"http.host":     req.Host,
	})
	h.handler.ServeHTTP(logger, req)
	dur := time.Now().UTC().Sub(t)
	sDur := float64(dur.Nanoseconds()) / (float64(time.Second) / float64(time.Nanosecond))
	logger.AddContext(logrus.Fields{
		"http.status": logger.Status(),
		"http.bytes":  logger.Size(),
		"http.dur":    sDur,
	})
}

// StructuredHandler returns an opinionated structuredHandler using the standard logger
// and setting a context with the fields:
// 	component = request.handler
func LogContextHandler(h http.Handler) http.Handler {
	return logContextHandler{h}
}
