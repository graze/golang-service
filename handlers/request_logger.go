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
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aristanetworks/goarista/monotime"
)

// LogServeHTTP creates a LoggingResponseWriter from `w` if applicable and calls `caller` with the request status, size,
// time and duration
func LogServeHTTP(w http.ResponseWriter, req *http.Request, handler http.Handler, caller func(w LoggingResponseWriter, req *http.Request, url url.URL, ts time.Time, dur time.Duration, status, size int)) {
	t := time.Now().UTC()
	mt := monotime.Now()
	logger := MakeLogger(w)
	url := *req.URL
	handler.ServeHTTP(logger, req)
	dur := time.Duration(monotime.Now() - mt)
	caller(logger, req, url, t, dur, logger.Status(), logger.Size())
}

// MakeLogger creates a LoggingResponseWriter from a http.ResponseWriter
//
// The loggingResponsWriter adds status field and the size of the response to the LoggingResponseWriter
func MakeLogger(w http.ResponseWriter) LoggingResponseWriter {
	if logResponse, ok := w.(LoggingResponseWriter); ok {
		return logResponse
	}

	var logger LoggingResponseWriter = &responseLogger{w: w}
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

// LoggingResponseWriter wraps a `http.ResponseWriter` `http.Flusher` and stores a logging context and status/size info
type LoggingResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	Status() int
	Size() int
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
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

// parseURI takes a request and url and returns a RFC7540 compliant uri
func parseURI(req *http.Request, url url.URL) (uri string) {
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

// getUserIP takes a request and extracts the users ip, using `X-Forwarded-For` and `X-Real-Ip` headers
// then the `req.RemoteAddr`
func getUserIP(req *http.Request) (net.IP, error) {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		if req.Header.Get(h) != "" {
			for _, ip := range strings.Split(req.Header.Get(h), ",") {
				// header can contain spaces too, strip those out.
				userIP := net.ParseIP(strings.Replace(ip, " ", "", -1))
				if userIP == nil {
					return nil, fmt.Errorf("getUserIP: %q is not a valid IP", ip)
				}
				return userIP, nil
			}
		}
	}
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("getUserIP: %q is not a valid IP:Port", req.RemoteAddr)
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("getUserIP: %q is not a valid IP:port", req.RemoteAddr)
	}
	return userIP, nil
}
