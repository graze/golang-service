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
	"fmt"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/graze/golang-service/logging"
)

type structuredHandler struct {
	logger  log.FieldLogger
	handler http.Handler
}

// ServeHTTP does the actual handling of HTTP requests by wrapping the request in a logger
func (h structuredHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	LogServeHTTP(w, req, h.handler, h.writeLog)
}

// writeLog writes a log entry to structuredHandler's logger
func (h structuredHandler) writeLog(req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int) {
	writeStructuredLog(h.logger, req, url, ts, dur, status, size)
}

// writeStructuredLog writes a log entry for req to logger in a structured format for json/logfmt
// ts is the timestamp with wich the entry should be logged
// dur is the time taken by the server to generate the response
// status and size are used to provide response HTTP status and size
func writeStructuredLog(logger log.FieldLogger, req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int) {
	sDur := float64(dur.Nanoseconds()) / (float64(time.Second) / float64(time.Nanosecond))
	uri := parseUri(req, url)

	logger.WithFields(log.Fields{
		"tag":      "request_handled",
		"method":   req.Method,
		"protocol": req.Proto,
		"uri":      uri,
		"path":     uriPath(req, url),
		"host":     req.Host,
		"status":   status,
		"bytes":    size,
		"dur":      sDur,
		"ts":       ts.Format(time.RFC3339Nano),
		"ref":      req.Referer(),
		"user":     req.Header.Get("X-Forwarded-For"),
	}).Info(fmt.Sprintf("%s %s %s", req.Method, uri, req.Proto))
}

// StructuredHandler return a http.Handler that wraps h and logs request to out in
// a structured format that can be outputted in json or logfmt
//
// Example:
//
//  r := mux.NewRouter()
//  r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//  	w.Write([]byte("This is a catch-all route"))
//  })
//  context := log.New().WithFields(log.Fields{
//	  "application": "service"
//	})
//  loggedRouter := logging.StructuredHandler(context, r)
//  http.ListenAndServe(":1123", loggedRouter)
//
func StructuredLogHandler(logger log.FieldLogger, h http.Handler) http.Handler {
	return structuredHandler{logger, h}
}

// StructuredHandler returns an opinionated structuredHandler using the standard logger
// and setting a context with the fields:
// 	component = request.handler
func StructuredHandler(h http.Handler) http.Handler {
	context := logging.GetLogger().WithFields(log.Fields{
		"component": "request.handler",
	})
	return structuredHandler{context, h}
}
