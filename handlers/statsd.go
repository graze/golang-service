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
	"strconv"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/graze/golang-service/logging"
)

type statsdHandler struct {
	statsd  *statsd.Client
	handler http.Handler
}

// ServeHTTP does the actual handling of HTTP requests by wrapping the request in a logger
func (h statsdHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	LogServeHTTP(w, req, h.handler, h.writeLog)
}

// writeLog writes the log do the statsd client from a statsdHandler
func (h statsdHandler) writeLog(w LoggingResponseWriter, req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int) {
	writeStatsdLog(h.statsd, req, url, ts, dur, status, size)
}

// writeStatsdLog send the response time and a counter for each request to statsd
func writeStatsdLog(w *statsd.Client, req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int) {
	uri := uriPath(req, url)

	tags := []string{
		"endpoint:" + uri,
		"statusCode:" + strconv.Itoa(status),
		"method:" + req.Method,
		"protocol:" + req.Proto,
	}

	msDur := float64(dur.Nanoseconds() / (int64(time.Millisecond) / int64(time.Nanosecond)))

	w.TimeInMilliseconds("request.response_time", msDur, tags, 1)
	w.Incr("request.count", tags, 1)
}

// StatsdIoHandler returns a http.Handler that wraps h and logs request to statsd
//
// Example:
//
//  r := mux.NewRouter()
//  r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//  	w.Write([]byte("This is a catch-all route"))
//  })
//  c, err := statsd.New("127.0.0.1:8125")
//  loggedRouter := logging.StatsdHandler(c, r)
//  http.ListenAndServe(":1123", loggedRouter)
//
func StatsdIoHandler(out *statsd.Client, h http.Handler) http.Handler {
	return statsdHandler{out, h}
}

// statsdHandler returns a logging.StatsdHandler to write request and response informtion to statsd
func StatsdHandler(h http.Handler) http.Handler {
	client, err := logging.GetStatsdFromEnv()
	if err != nil {
		panic(err)
	}
	return StatsdIoHandler(client, h)
}
