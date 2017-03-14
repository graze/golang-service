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

	"github.com/graze/golang-service/log"
	uuid "github.com/satori/go.uuid"
)

// logContextHandler contains a local logger context and the handler
type logContextHandler struct {
	logger  log.FieldLogger
	handler http.Handler
}

// ServeHTTP modifies the context of the request by adding a local logger context
func (h logContextHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := *req.URL
	ip := ""
	if userIP, err := getUserIP(req); err == nil {
		ip = userIP.String()
	}
	ctx := h.logger.Ctx(req.Context()).With(log.KV{
		"transaction":     uuid.NewV4().String(),
		"http.method":     req.Method,
		"http.protocol":   req.Proto,
		"http.uri":        parseURI(req, url),
		"http.path":       uriPath(req, url),
		"http.host":       req.Host,
		"http.user":       ip,
		"http.ref":        req.Referer(),
		"http.user-agent": req.Header.Get("User-Agent"),
	}).NewContext(req.Context())
	h.handler.ServeHTTP(w, req.WithContext(ctx))
}

// LoggingContextHandler returns a handler that adds `http` and `transaction` items into the provided logging context
//
// It adds the following fields to the `LoggingResponseWriter` log context:
//  http.method     - GET/POST/...
//  http.protocol   - HTTP/1.1
//  http.uri        - /path?with=query
//  http.path       - /path
//  http.host       - localhost:80
//	http.user		- 192.168.0.1 - ip address of the user
//	http.ref		- http://google.com - referrer
//	http.user-agent - The user agent of the user
//  transaction     - unique uuid4 for this request
func LoggingContextHandler(logger log.FieldLogger, h http.Handler) http.Handler {
	return logContextHandler{logger.With(log.KV{}), h}
}
