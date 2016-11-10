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
	"net/url"
	"time"

	"github.com/graze/golang-service/log"
)

type structuredHandler struct {
	logger  log.Context
	handler http.Handler
}

// ServeHTTP does the actual handling of HTTP requests by wrapping the request in a logger
func (h structuredHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	LogServeHTTP(w, req, h.handler, h.writeLog)
}

// writeLog writes a log entry to structuredHandler's logger
func (h structuredHandler) writeLog(w LoggingResponseWriter, req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int) {
	writeStructuredLog(w, h.logger.Ctx(req.Context()), req, url, ts, dur, status, size)
}

// writeStructuredLog writes a log entry for req to logger in a structured format for json/logfmt
// ts is the timestamp with wich the entry should be logged
// dur is the time taken by the server to generate the response
// status and size are used to provide response HTTP status and size
func writeStructuredLog(w LoggingResponseWriter, logger log.Context, req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int) {
	uri := parseURI(req, url)
	ip := ""
	if userIP, err := getUserIP(req); err == nil {
		ip = userIP.String()
	}

	logger.With(log.KV{
		"tag":             "request_handled",
		"http.method":     req.Method,
		"http.protocol":   req.Proto,
		"http.uri":        uri,
		"http.path":       uriPath(req, url),
		"http.host":       req.Host,
		"http.status":     status,
		"http.bytes":      size,
		"http.user":       ip,
		"http.ref":        req.Referer(),
		"http.user-agent": req.Header.Get("User-Agent"),
		"dur":             dur.Seconds(),
		"ts":              ts.Format(time.RFC3339Nano),
	}).Infof("%s %s %s", req.Method, uri, req.Proto)
}

// StructuredLogHandler returns a http.Handler that wraps h and logs request to out in
// a structured format that can be outputted in json or logfmt
//
// Example:
//
//  r := mux.NewRouter()
//  r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//  	w.Write([]byte("This is a catch-all route"))
//  })
//  context := log.New()
//	context.SetFormatter(&logrus.JSONFormatter{})
//	context.Add(log.KV{"module":"request.handler"})
//  loggedRouter := handlers.StructuredLogHandler(
//		log.With(log.KV{"module":"request.handler"})
//		, r)
//  http.ListenAndServe(":1123", loggedRouter)
func StructuredLogHandler(logger log.Context, h http.Handler) http.Handler {
	return structuredHandler{logger, h}
}

// StructuredHandler returns an opinionated structuredHandler using the standard logger
// and setting a context with the fields:
// 	component = request.handler
//
// Usage:
//  r := mux.NewRouter()
//  r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//  	w.Write([]byte("This is a catch-all route"))
//  })
//  loggedRouter := handlers.StructuredHandler(r)
//  http.ListenAndServe(":1123", loggedRouter)
//
func StructuredHandler(h http.Handler) http.Handler {
	logger := log.With(log.KV{
		"module": "request.handler",
	})
	return structuredHandler{logger, h}
}
