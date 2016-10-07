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
    "io"
    "net/http"
	"net/url"
    "time"
    "strconv"
    "fmt"
)

type healthdHandler struct {
    writer  io.Writer
    handler http.Handler
}

// ServeHTTP does the actual handling of HTTP requests by wrapping the request in a logger
func (h healthdHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    t := time.Now().UTC()
	logger := MakeLogger(w)
	url := *req.URL
	h.handler.ServeHTTP(logger, req)
    dur := time.Now().UTC().Sub(t)
	writeHealthdLog(h.writer, req, url, t, dur, logger.Status(), logger.Size())
}

// writeHealthdLog writes a log entry for req to w in healthd format.
// ts is the timestamp with wich the entry should be logged
// dur is the time taken by the server to generate the response
// status and size are used to provide response HTTP status and size
//
// The format of the file is:
// <unix_timestamp.ms>"<path>"<status>"<request_time>"<upstream_time>"<X-Forwarded-For header>
func writeHealthdLog(w io.Writer, req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int) {
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

    msDur := float64(dur.Nanoseconds()) / (float64(time.Second)/float64(time.Nanosecond))
    str := fmt.Sprintf(`%.3f%s%d"%.3f"%.3f"%s` + "\n",
        float64(ts.UnixNano()) / (float64(time.Second)/float64(time.Nanosecond)),
        strconv.Quote(uri),
        status,
        msDur,
        msDur,
        req.Header.Get("X-Forwarded-For"))
    io.WriteString(w, str)
}

// HealthdHandler return a http.Handler that wraps h and logs request to out in
// nginx Healthd format
//
// see http://docs.aws.amazon.com/elasticbeanstalk/latest/dg/health-enhanced-serverlogs.html
//
// Example:
//
//  r := mux.NewRouter()
//  r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//  	w.Write([]byte("This is a catch-all route"))
//  })
//  loggedRouter := logging.HealthdHandler(os.Stdout, r)
//  http.ListenAndServe(":1123", loggedRouter)
//
func HealthdHandler(out io.Writer, h http.Handler) http.Handler {
    return healthdHandler{out, h}
}
