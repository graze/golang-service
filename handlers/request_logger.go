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
	"bufio"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
	log "github.com/graze/golang-service/logging"
	"github.com/twinj/uuid"
)

func LogServeHTTP(w http.ResponseWriter, req *http.Request, handler http.Handler, caller func(w LoggingResponseWriter, req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int)) {
	t := time.Now().UTC()
	logger := MakeLogger(w)
	url := *req.URL
	handler.ServeHTTP(logger, req)
	dur := time.Now().UTC().Sub(t)
	caller(logger, req, url, t, dur, logger.Status(), logger.Size())
}

// MakeLogger creates a LoggingResponseWriter from a http.ResponseWriter
//
// The loggingResponsWriter adds status field and the size of the response to the LoggingResponseWriter
func MakeLogger(w http.ResponseWriter) LoggingResponseWriter {
	if _, ok := w.(LoggingResponseWriter); ok {
		return w.(LoggingResponseWriter)
	}

	context := log.WithFields(logrus.Fields{
		"transaction": uuid.NewV4(),
	})
	var logger LoggingResponseWriter = &responseLogger{w: w, Context: context}
	if _, ok := w.(http.Hijacker); ok {
		logger = &hijackLogger{responseLogger{w: w}}
	}
	h, ok1 := logger.(http.Hijacker)
	c, ok2 := w.(http.CloseNotifier)
	if ok1 && ok2 {
		return hijackCloseNotifier{logger, h, c}
	}
	if ok2 {
		return &closeNotifyWriter{logger, c}
	}
	return logger
}

type LoggingResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	Status() int
	Size() int
	AddContext(fields logrus.Fields)
	GetContext() log.LogContext
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size
type responseLogger struct {
	w       http.ResponseWriter
	Context *logrus.Entry
	status  int
	size    int
}

// AddContext appends items to the current logging context for this http request
func (l *responseLogger) AddContext(fields logrus.Fields) {
	l.Context = l.Context.WithFields(fields)
}

// GetContext gets the current logging context for this http request
func (l *responseLogger) GetContext() log.LogContext {
	return l.Context
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.status = http.StatusOK
	}
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func (l *responseLogger) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

type hijackLogger struct {
	responseLogger
}

func (l *hijackLogger) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h := l.responseLogger.w.(http.Hijacker)
	conn, rw, err := h.Hijack()
	if err == nil && l.responseLogger.status == 0 {
		// The status will be StatusSwitchingProtocols if there was no error and
		// WriteHeader has not been called yet
		l.responseLogger.status = http.StatusSwitchingProtocols
	}
	return conn, rw, err
}

type closeNotifyWriter struct {
	LoggingResponseWriter
	http.CloseNotifier
}

type hijackCloseNotifier struct {
	LoggingResponseWriter
	http.Hijacker
	http.CloseNotifier
}

// parseUri takes a request and url and returns a RFC7540 compliant uri
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

// uriPath extracts the path from a request uri
func uriPath(req *http.Request, url url.URL) (uri string) {
	uri = url.EscapedPath()

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
	return
}
