// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// license: https://github.com/graze/golang-service/blob/master/LICENSE
// link:    https://github.com/graze/golang-service

package logging

import (
    "github.com/DataDog/datadog-go/statsd"
    "net/http"
    "net/url"
    "time"
    "strconv"
)

type statsdHandler struct {
    statsd  *statsd.Client
    handler http.Handler
}

// ServeHTTP does the actual handling of HTTP requests by wrapping the request in a logger
func (h statsdHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    t := time.Now().UTC()
	logger := MakeLogger(w)
	url := *req.URL
	h.handler.ServeHTTP(logger, req)
    dur := time.Now().UTC().Sub(t)
	writeStatsdLog(h.statsd, req, url, t, dur, logger.Status(), logger.Size())
}

// writeStatsdLog send the response time and a counter for each request to statsd
func writeStatsdLog(w *statsd.Client, req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int) {
    uri := url.EscapedPath()

    // Requests using the CONNECT method over HTTP/2.0 must use
    // the authority field (aka r.Host) to identify the target.
    // Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
    if req.ProtoMajor == 2 && req.Method == "CONNECT" {
        uri = req.Host
    }
    if uri == "" {
        parsed, err := url.Parse(req.RequestURI)
        if err != nil {
            uri = "unknown"
        } else {
            uri = parsed.EscapedPath()
        }
    }
    if uri == "" {
        uri = "/"
    }

    tags := []string{
        "endpoint:" + uri,
        "statusCode:" + strconv.Itoa(status),
        "method:" + req.Method,
        "protocol:" + req.Proto,
    }

    msDur := float64(dur.Nanoseconds() / (int64(time.Millisecond)/int64(time.Nanosecond)))

    w.TimeInMilliseconds("request.response_time", msDur, tags, 1)
    w.Incr("request.count", tags, 1)
}

// StatsdHandler returns a http.Handler that wraps h and logs request to statsd
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
func StatsdHandler(out *statsd.Client, h http.Handler) http.Handler {
    return statsdHandler{out, h}
}
