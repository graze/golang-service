// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// @license https://github.com/graze/golang-service-logging/blob/master/LICENSE
// @link    https://github.com/graze/golang-service-logging
package logging


import (
    "github.com/DataDog/datadog-go/statsd"
    "net/http"
    "net/url"
    "time"
    "strconv"
    "strings"
)

type statsdHandler struct {
    statsd  *statsd.Client
    handler http.Handler
}

func (h statsdHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    t := time.Now()
	logger := MakeLogger(w)
	url := *req.URL
	h.handler.ServeHTTP(logger, req)
    dur := time.Since(t)
	writeStatsdLog(h.statsd, req, url, t, dur, logger.Status(), logger.Size())
}

// writeStatsdLog send the response time and a counter for each request to statsd
func writeStatsdLog(w *statsd.Client, req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int) {
    uri := req.RequestURI

    // Requests using the CONNECT method over HTTP/2.0 must use
    // the authority field (aka r.Host) to identify the target.
    // Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
    if req.ProtoMajor == 2 && req.Method == "CONNECT" {
        uri = req.Host
    }
    if uri == "" {
        uri = url.RequestURI()
    }

    tags := []string{
        strings.Join([]string{"endpoint", uri}, ":"),
        strings.Join([]string{"statusCode", strconv.Itoa(status)}, ":"),
        strings.Join([]string{"method", req.Method}, ":"),
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
