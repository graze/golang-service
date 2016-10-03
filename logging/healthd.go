// This file is part of graze/golang-service-logging
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
    "io"
    "net/http"
	"net/url"
    "time"
    "strconv"
)

type healthdHandler struct {
    writer  io.Writer
    handler http.Handler
}

func (h healthdHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    t := time.Now()
	logger := MakeLogger(w)
	url := *req.URL
	h.handler.ServeHTTP(logger, req)
    dur := time.Since(t)
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

    buf := make([]byte, 0)
    buf = append(buf, strconv.FormatFloat(float64(ts.UnixNano()) / (float64(time.Second)/float64(time.Nanosecond)), 'f', 3, 64)...)
    buf = append(buf, `"`...)
    buf = AppendQuoted(buf, uri)
    buf = append(buf, `"`...)
    buf = append(buf, strconv.Itoa(status)...)
    buf = append(buf, `"`...)

    msDur := float64(dur.Nanoseconds()) / (float64(time.Second)/float64(time.Nanosecond))
    buf = append(buf, strconv.FormatFloat(msDur, 'f', 3, 64)...)
    buf = append(buf, `"`...)
    buf = append(buf, strconv.FormatFloat(msDur, 'f', 3, 64)...)
    buf = append(buf, `"`...)
    buf = append(buf, req.Header.Get("X-Forwarded-For")...)
    buf = append(buf, '\n')
    w.Write(buf)
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
