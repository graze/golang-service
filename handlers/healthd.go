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
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type healthdHandler struct {
	writer  io.Writer
	handler http.Handler
}

// ServeHTTP does the actual handling of HTTP requests by wrapping the request in a logger
func (h healthdHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	LogServeHTTP(w, req, h.handler, h.writeLog)
}

// writeLog writes a log entry to healthdHandler's writer
func (h healthdHandler) writeLog(req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int) {
	writeHealthdLog(h.writer, req, url, ts, dur, status, size)
}

func parseUri(req *http.Request, url url.URL) (uri string) {
	uri = req.RequestURI

	// Requests using the CONNECT method over HTTP/2.0 must use
	// the authority field (aka r.Host) to identify the target.
	// Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
	if req.ProtoMajor == 2 && req.Method == "CONNECT" {
		uri = req.Host
	}
	if uri == "" {
		uri = url.RequestURI()
	}
	return
}

// writeHealthdLog writes a log entry for req to w in healthd format.
// ts is the timestamp with wich the entry should be logged
// dur is the time taken by the server to generate the response
// status and size are used to provide response HTTP status and size
//
// The format of the file is:
// <unix_timestamp.ms>"<path>"<status>"<request_time>"<upstream_time>"<X-Forwarded-For header>
func writeHealthdLog(w io.Writer, req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int) {
	uri := parseUri(req, url)
	msDur := float64(dur.Nanoseconds()) / (float64(time.Second) / float64(time.Nanosecond))
	str := fmt.Sprintf(`%.3f%s%d"%.3f"%.3f"%s`+"\n",
		float64(ts.UnixNano())/(float64(time.Second)/float64(time.Nanosecond)),
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
func HealthdIoHandler(out io.Writer, h http.Handler) http.Handler {
	return healthdHandler{out, h}
}

type healthdFileHandler struct {
	path      string
	base      http.Handler
	handler   http.Handler
	timestamp string
}

// ServerHTTP for the healthdFileHandler automatically rotates the log files based on the hour
func (h healthdFileHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	now := time.Now().UTC().Format("2006-01-02-15")
	if h.timestamp != now || h.base == nil {
		h.timestamp = now
		file := h.path + "application.log." + now
		logFile, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()
		h.base = HealthdIoHandler(logFile, h.handler)
	}
	h.base.ServeHTTP(w, req)
}

// healthdHandler returns a logging.HealthdHandler and outputs healthd formatted log files to the appropriate path
func HealthdHandler(h http.Handler) http.Handler {
	path := "/var/log/nginx/healthd/"
	err := os.MkdirAll(path, 0666)
	if err != nil {
		panic(err)
	}

	return healthdFileHandler{path, nil, h, ""}
}
